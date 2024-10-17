package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

func init() {
	// Load environment variables
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Printf("Error loading .env.local file: %v", err)
	}

	// Redis initialization
	redisOptions, err := redis.ParseURL(os.Getenv("REDIS_UPSTASH_ADDR"))
	if err != nil {
		log.Fatalf("Error parsing Redis URL: %v", err)
	}
	redisClient = redis.NewClient(redisOptions)

	// Test Redis connection
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}

	log.Println("Connected to Redis successfully.")
}

func SaveDriverLocationToRedis(driverID int, lat, lng float64) {
	key := fmt.Sprintf("driver_location:%d", driverID)
	loc := Location{
		Lat: lat,
		Lng: lng,
	}
	locBytes, _ := json.Marshal(loc)
	err := redisClient.Set(ctx, key, locBytes, time.Minute*10).Err()
	if err != nil {
		log.Printf("Failed to save driver location to Redis: %v", err)
	} else {
		log.Printf("Driver location saved to Redis for driver %d", driverID)
	}
}
