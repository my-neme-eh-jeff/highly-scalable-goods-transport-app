package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func main() {
    // Initialize MongoDB and Redis clients
    initMongoDB()
    initRedis()

    http.HandleFunc("/ws/driver/update-location", DriverLocationWebSocket)
    log.Println("Listening on :8083")
    http.ListenAndServe(":8083", nil)
}

func initMongoDB() {
    var err error
    mongoClient, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatalf("Failed to create MongoDB client: %v", err)
    }
    err = mongoClient.Connect(ctx)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    log.Println("Connected to MongoDB.")
}

func initRedis() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    err := redisClient.Ping(ctx).Err()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    log.Println("Connected to Redis.")
}


