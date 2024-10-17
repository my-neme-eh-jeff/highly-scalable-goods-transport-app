package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	dbPool      *pgxpool.Pool
)

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Printf("Error loading .env.local file: %v", err)
	}

	opt, err := redis.ParseURL(os.Getenv("REDIS_UPSTASH_ADDR"))
	if err != nil {
		log.Printf("Error parsing Redis URL: %v", err)
	}

	opt.IdleTimeout = time.Minute * 5

	redisClient = redis.NewClient(opt)

	config, err := pgxpool.ParseConfig(os.Getenv("NEON_DB_URL"))
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Test connections
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Printf("Redis connection error: %v", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Printf("Database connection error: %v", err)
	}

	_, err = dbPool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			pickup_location GEOGRAPHY(POINT),
			dropoff_location GEOGRAPHY(POINT),
			distance_km DECIMAL(10,2),
			fare_amount DECIMAL(10,2),
			status VARCHAR(20),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("Error creating table: %v", err)
	}
}

func GetFare(w http.ResponseWriter, r *http.Request) {
	var req FareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate distance
	distance := CalculateDistance(req.PickupLocation.Lat, req.PickupLocation.Lng,
		req.DropoffLocation.Lat, req.DropoffLocation.Lng)

	fare := distance * 50 // 1 km = 50 rupees

	// Check for surge pricing
	surgeMultiplier := CheckSurgePricing(req.PickupLocation)
	fare *= surgeMultiplier

	// Save fare request in Redis payment cache
	SaveFareRequestInRedis(req.UserID, req.PickupLocation)

	// Create payment record
	payment := Payment{
		UserID:          req.UserID,
		PickupLocation:  fmt.Sprintf("POINT(%f %f)", req.PickupLocation.Lng, req.PickupLocation.Lat),
		DropoffLocation: fmt.Sprintf("POINT(%f %f)", req.DropoffLocation.Lng, req.DropoffLocation.Lat),
		DistanceKM:      distance,
		FareAmount:      fare,
		Status:          "PENDING",
	}

	// Insert into database using connection from pool
	var paymentID int
	err := dbPool.QueryRow(ctx, `
		INSERT INTO payments (user_id, pickup_location, dropoff_location, distance_km, fare_amount, status)
		VALUES ($1, ST_GeographyFromText($2), ST_GeographyFromText($3), $4, $5, $6)
		RETURNING id
	`,
		payment.UserID,
		fmt.Sprintf("SRID=4326;%s", payment.PickupLocation),
		fmt.Sprintf("SRID=4326;%s", payment.DropoffLocation),
		payment.DistanceKM,
		payment.FareAmount,
		payment.Status,
	).Scan(&paymentID)

	if err != nil {
		log.Printf("Error saving payment: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	resp := map[string]interface{}{
		"fare_amount": fare,
		"distance_km": distance,
		"status":      "PENDING",
		"payment_id":  paymentID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
