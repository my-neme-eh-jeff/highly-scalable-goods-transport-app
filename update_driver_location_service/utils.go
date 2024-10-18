package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/geo/s2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx         = context.Background()
	mongoClient *mongo.Client
	redisClient *redis.Client
)

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Printf("Error loading .env.local file: %v", err)
	}
	initMongoDB()
	initRedis()
}

func initMongoDB() {
	var err error
	mongoURI := os.Getenv("MONGODB_URI")
	// fmt.Println("MONGODB_URI: ", mongoURI)
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	log.Println("Connected to MongoDB.")
}

func initRedis() {
	redisURL := os.Getenv("REDIS_UPSTASH_ADDR")
	if redisURL == "" {
		log.Fatal("UPSTASH_REDIS_URL environment variable not set")
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	redisClient = redis.NewClient(options)

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Upstash Redis.")
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
	} else {
		log.Printf("Driver location saved to Redis with key %s", key)
	}

	driverLatLng := s2.LatLngFromDegrees(location.Lat, location.Lng)
	driverCellID := s2.CellIDFromLatLng(driverLatLng).Parent(15)
	cellKey := fmt.Sprintf("drivers_in_cell:%d", driverCellID)

	// Add driver ID to the set for the cell
	redisClient.SAdd(ctx, cellKey, driverID).Err()

	// Set expiration for the cell key
	err = redisClient.Expire(ctx, cellKey, time.Minute*10).Err()
	if err != nil {
		log.Printf("Failed to set expiration for cell key: %v", err)
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
	} else {
		log.Printf("Location saved to MongoDB for booking %s", bookingID)
	}
}
