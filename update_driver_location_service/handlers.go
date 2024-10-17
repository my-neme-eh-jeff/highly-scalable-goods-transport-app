package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	upgrader    = websocket.Upgrader{}
	ctx         = context.Background()
	mongoClient *mongo.Client
	redisClient *redis.Client
)

func DriverLocationWebSocket(w http.ResponseWriter, r *http.Request) {
	driverIDStr := r.URL.Query().Get("driver_id")
	bookingID := r.URL.Query().Get("booking_id")
	if driverIDStr == "" || bookingID == "" {
		http.Error(w, "Missing driver_id or booking_id parameter", http.StatusBadRequest)
		return
	}

	driverID, err := strconv.Atoi(driverIDStr)
	if err != nil {
		http.Error(w, "Invalid driver_id parameter", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		var locationUpdate struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}
		err := conn.ReadJSON(&locationUpdate)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Save to Redis
		SaveDriverLocationToRedis(driverID, Location{Lat: locationUpdate.Lat, Lng: locationUpdate.Lng})

		// Save to MongoDB for tracking
		SaveLocationToMongoDB(bookingID, locationUpdate.Lat, locationUpdate.Lng)
	}
}
