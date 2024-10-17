package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	ctx         = context.Background()
	dbPool      *pgxpool.Pool
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

	err = ensureLocationsCollectionExists()
	if err != nil {
		log.Fatalf("Error ensuring MongoDB 'locations' collection: %v", err)
	}

	rabbitMQURL := os.Getenv("CLOUDAMQP_URL")
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to RabbitMQ successfully.")
}

// EnsureLocationsCollectionExists checks if the locations collection exists in MongoDB
func ensureLocationsCollectionExists() error {
	collection := mongoClient.Database("transport").Collection("locations")
	indexes := collection.Indexes()

	// Create an index on booking_id and timestamp fields
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

func SaveLocationToMongoDB(bookingID string, location string) {
	collection := mongoClient.Database("transport").Collection("locations")
	_, err := collection.InsertOne(ctx, map[string]interface{}{
		"booking_id": bookingID,
		"location":   location,
		"timestamp":  time.Now(),
	})
	if err != nil {
		log.Printf("Error saving location to MongoDB: %v", err)
	} else {
		log.Printf("Location saved to MongoDB for booking %s", bookingID)
	}
}

func GetLocationsFromMongoDB(bookingID string) []string {
	collection := mongoClient.Database("transport").Collection("locations")
	cursor, err := collection.Find(ctx, map[string]interface{}{
		"booking_id": bookingID,
	})
	if err != nil {
		log.Printf("Error retrieving locations from MongoDB: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	var locations []string
	for cursor.Next(ctx) {
		var result map[string]interface{}
		cursor.Decode(&result)
		locations = append(locations, result["location"].(string))
	}
	return locations
}

func GetLatestLocation(bookingID string) string {
	collection := mongoClient.Database("transport").Collection("locations")
	var result map[string]interface{}
	err := collection.FindOne(ctx, map[string]interface{}{
		"booking_id": bookingID,
	}, options.FindOne().SetSort(map[string]int{"timestamp": -1})).Decode(&result)
	if err != nil {
		log.Printf("Error getting latest location from MongoDB for booking %s: %v", bookingID, err)
		return ""
	}
	return result["location"].(string)
}

func GetMongoChangeStream(bookingID string) *mongo.ChangeStream {
	collection := mongoClient.Database("transport").Collection("locations")

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"fullDocument.booking_id", bookingID},
		}}},
	}

	opts := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	changeStream, err := collection.Watch(ctx, pipeline, opts)
	if err != nil {
		log.Printf("Error creating change stream: %v", err)
	}
	return changeStream
}
