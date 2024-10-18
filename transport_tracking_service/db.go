package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
    mongoClient *mongo.Client
    ctx         = context.Background()
)

func init() {
    err := godotenv.Load(".env.local")
    if err != nil {
        log.Printf("Error loading .env.local file: %v", err)
    }

    mongoURI := os.Getenv("MONGODB_URI")
    clientOptions := options.Client().ApplyURI(mongoURI)
    mongoClient, err = mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatalf("Error connecting to MongoDB: %v", err)
    }

    err = mongoClient.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("MongoDB connection error: %v", err)
    }

    log.Println("Connected to MongoDB successfully.")

    // Ensure indexes
    err = ensureLocationsCollectionExists()
    if err != nil {
        log.Fatalf("Error ensuring MongoDB 'locations' collection: %v", err)
    }
}

func ensureLocationsCollectionExists() error {
    collection := mongoClient.Database("transport").Collection("locations")
    indexes := collection.Indexes()

    indexModel := mongo.IndexModel{
        Keys: bson.D{
            {Key: "booking_id", Value: 1},
            {Key: "timestamp", Value: 1},
        },
        Options: options.Index().SetUnique(false),
    }

    _, err := indexes.CreateOne(ctx, indexModel)
    if err != nil {
        return fmt.Errorf("error creating index for 'locations' collection: %v", err)
    }
    return nil
}
