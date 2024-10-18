package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func BookTransport(w http.ResponseWriter, r *http.Request) {
	var req BookingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create booking record in Postgres
	bookingID, err := SaveBooking(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.BookingID = bookingID // Set the booking ID in the request

	// Push to Kafka
	err = PushToKafka("user_bookings", req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	resp := map[string]interface{}{
		"booking_id": bookingID,
		"status":     "REQUESTED",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func DriverRespondBooking(w http.ResponseWriter, r *http.Request) {
	var req DriverResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Response == "ACCEPT" {
		// Update booking status to 'STARTED'
		err := UpdateBookingStatus(req.BookingID, "STARTED")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Push event to RabbitMQ or perform other actions as needed
		// For example, notify the user that the driver has accepted the booking
	} else if req.Response == "REJECT" {
		// Release driver lock in Redis
		ReleaseDriverLock(req.DriverID)
		// Update booking status to 'PENDING' or 'REJECTED'
		err := UpdateBookingStatus(req.BookingID, "PENDING")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Optionally, you can implement logic to try assigning another driver
	}

	// Respond to driver
	resp := map[string]string{
		"status": "OK",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

    userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	bookings, err := GetBookingsByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//print the bookings
	fmt.Println(bookings)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func DriverCompleteRide(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DriverID  int `json:"driver_id"`
		BookingID int `json:"booking_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update booking status to 'COMPLETED'
	err = UpdateBookingStatus(req.BookingID, "COMPLETED")
	if err != nil {
		http.Error(w, "Failed to update booking status", http.StatusInternalServerError)
		return
	}

	// Release driver lock
	ReleaseDriverLock(req.DriverID)

	// Respond to the driver
	resp := map[string]string{
		"status": "COMPLETED",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
