package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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

func SaveDriverLocationToRedis(driverID int, location Location) {
    key := fmt.Sprintf("driver_location:%d", driverID)
    locationJSON, err := json.Marshal(location)
    if err != nil {
        log.Printf("Error marshaling location: %v", err)
        return
    }

    err = redisClient.Set(ctx, key, locationJSON, 0).Err()
    if err != nil {
        log.Printf("Error saving driver location to Redis: %v", err)
    }
}

func SaveLocationToMongoDB(bookingID string, lat, lng float64) {
    collection := mongoClient.Database("transport").Collection("locations")
    _, err := collection.InsertOne(ctx, LocationRecord{
        BookingID: bookingID,
        Location:  Location{Lat: lat, Lng: lng},
        Timestamp: time.Now(),
    })
    if err != nil {
        log.Printf("Error saving location to MongoDB: %v", err)
    }
}
