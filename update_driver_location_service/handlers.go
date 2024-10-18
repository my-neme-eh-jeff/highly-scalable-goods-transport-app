package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development; secure this in production
	},
}

func DriverLocationWebSocket(w http.ResponseWriter, r *http.Request) {
	driverIDStr := r.URL.Query().Get("driver_id")
	bookingIDStr := r.URL.Query().Get("booking_id")
	log.Printf("Driver %s connected for location updates.", driverIDStr)
	if driverIDStr == "" {
		http.Error(w, "Missing driver_id parameter", http.StatusBadRequest)
		return
	}

	driverID, err := strconv.Atoi(driverIDStr)
	if err != nil {
		http.Error(w, "Invalid driver_id parameter", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Driver %d connected for location updates.", driverID)

	for {
		var locationUpdate struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		}
		err := conn.ReadJSON(&locationUpdate)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		if bookingIDStr != "" {
			// Driver is in a ride, save to MongoDB
			SaveLocationToMongoDB(bookingIDStr, locationUpdate.Lat, locationUpdate.Lng)
		} else {
			// Driver is not in a ride, save to Redis
			driverLocation := Location{
				Lat: locationUpdate.Lat,
				Lng: locationUpdate.Lng,
			}
			SaveDriverLocationToRedis(driverID, driverLocation)
		}
	}

	log.Printf("Driver %d disconnected from location updates.", driverID)
}
