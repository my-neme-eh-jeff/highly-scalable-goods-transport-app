package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/geo/s2"
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
	// Redis initialization
	redisOptions, err := redis.ParseURL(os.Getenv("REDIS_UPSTASH_ADDR"))
	if err != nil {
		log.Fatal("Error afasdasdasd %v", err)
		log.Fatalf("Error parsing Redis URL\n\n\n\n: %v", err)
	}
	redisClient = redis.NewClient(redisOptions)

	// Postgres connection pool
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("NEON_DB_URL"))
	if err != nil {
		log.Fatalf("Unable to parse Postgres config: %v", err)
	}
	dbPool, err = pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Test connections
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Fatalf("Postgres connection error: %v", err)
	}

	// Ensure the 'bookings' table exists
	_, err = dbPool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS bookings (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            driver_id INTEGER,
            pickup_location GEOGRAPHY(POINT),
            dropoff_location GEOGRAPHY(POINT),
            fare_amount DECIMAL(10,2),
            status VARCHAR(20),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatalf("Error creating 'bookings' table: %v", err)
	}

	log.Println("Connected to Redis and Postgres, and ensured 'bookings' table exists.")
}
func FindNearbyDrivers(pickup Location) ([]Driver, error) {
	pickupLatLng := s2.LatLngFromDegrees(pickup.Lat, pickup.Lng)
	pickupCellID := s2.CellIDFromLatLng(pickupLatLng).Parent(15)

	driverIDs, err := redisClient.SMembers(ctx, fmt.Sprintf("drivers_in_cell:%d", pickupCellID)).Result()
	if err != nil && err != redis.Nil {
		log.Printf("Redis SMembers error: %v", err)
		return nil, err
	}

	var nearbyDrivers []Driver
	for _, driverIDStr := range driverIDs {
		driverID, _ := strconv.Atoi(driverIDStr)
		loc, err := GetDriverLocation(driverID)
		if err != nil {
			continue
		}
		nearbyDrivers = append(nearbyDrivers, Driver{DriverID: driverID, Lat: loc.Lat, Lng: loc.Lng})
	}

	return nearbyDrivers, nil
}

func GetDriverLocation(driverID int) (Location, error) {
	key := fmt.Sprintf("driver_location:%d", driverID)
	driverLoc, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return Location{}, err
	}

	var loc Location
	err = json.Unmarshal([]byte(driverLoc), &loc)
	return loc, err
}

func AcquireDriverLock(driverID int) bool {
	lockKey := fmt.Sprintf("driver_lock:%d", driverID)
	success, err := redisClient.SetNX(ctx, lockKey, "locked", time.Minute).Result()
	if err != nil {
		log.Printf("Failed to acquire driver lock: %v", err)
		return false
	}
	return success
}

func ReleaseDriverLock(driverID int) {
	lockKey := fmt.Sprintf("driver_lock:%d", driverID)
	err := redisClient.Del(ctx, lockKey).Err()
	if err != nil {
		log.Printf("Failed to release driver lock: %v", err)
	}
}
