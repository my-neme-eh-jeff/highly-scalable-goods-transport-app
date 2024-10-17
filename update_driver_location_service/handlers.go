package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func DriverLocationWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		var locationUpdate struct {
			DriverID int     `json:"driver_id"`
			Lat      float64 `json:"lat"`
			Lng      float64 `json:"lng"`
		}
		err := conn.ReadJSON(&locationUpdate)
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		SaveDriverLocationToRedis(locationUpdate.DriverID, locationUpdate.Lat, locationUpdate.Lng)
	}
}
