package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func MatchDrivers(w http.ResponseWriter, r *http.Request) {
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
			NotifyDriver(driver.DriverID, booking)

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

// AssignDriver updates the booking with the assigned driver
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

// NotifyDriver sends the booking details to the driver via WebSocket
func NotifyDriver(driverID int, booking Booking) {
	// Establish WebSocket connection with the driver
	url := fmt.Sprintf("ws://localhost:8083/ws/driver/notify?driver_id=%d", driverID)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Printf("Error connecting to WebSocket for driver %d: %v", driverID, err)
		return
	}
	defer conn.Close()

	// Send booking details to the driver
	err = conn.WriteJSON(booking)
	if err != nil {
		log.Printf("Error sending booking details to driver %d: %v", driverID, err)
	}
}
