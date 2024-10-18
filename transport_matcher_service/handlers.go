package main

import (
	"encoding/json"
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

func MatchDrivers(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var booking Booking
		if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Find nearby drivers
		drivers, err := FindNearbyDrivers(booking.PickupLocation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Assign the first available driver
		for _, driver := range drivers {
			if AcquireDriverLock(driver.DriverID) {
				// Assign the driver
				err := AssignDriver(booking.BookingID, driver.DriverID)
				if err != nil {
					ReleaseDriverLock(driver.DriverID) // Release the lock if assignment fails
					continue
				}

				// Notify the driver via WebSocket
				NotifyDriver(hub, driver.DriverID, booking)

				// Respond to the client
				resp := map[string]interface{}{
					"booking_id": booking.BookingID,
					"driver_id":  driver.DriverID,
					"status":     "DRIVER_ASSIGNED",
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		// No driver found
		http.Error(w, "No drivers available", http.StatusServiceUnavailable)
	}
}

func AssignDriver(bookingID, driverID int) error {
	_, err := dbPool.Exec(ctx, `
        UPDATE bookings SET driver_id = $1, status = $2 WHERE id = $3
    `, driverID, "DRIVER_ASSIGNED", bookingID)
	if err != nil {
		log.Printf("Error assigning driver to booking %d: %v", bookingID, err)
		return err
	}
	return nil
}

func NotifyDriver(hub *Hub, driverID int, booking Booking) {
	hub.SendToDriver(driverID, booking)
}

func ServeDriverWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	driverIDStr := r.URL.Query().Get("driver_id")
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

	client := &Client{
		driverID: driverID,
		hub:      hub,
		conn:     conn,
		send:     make(chan interface{}, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
