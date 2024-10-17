package main

import (
	"encoding/json"
	"net/http"
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
		// Push event to RabbitMQ
		PushEventToRabbitMQ("booking_started", req.BookingID)
	} else {
		// Release driver lock in Redis
		ReleaseDriverLock(req.DriverID)
	}

	// Respond to driver
	resp := map[string]string{
		"status": "OK",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
