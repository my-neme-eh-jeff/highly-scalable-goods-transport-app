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

	req.BookingID = bookingID 

	// Push to Kafka
	err = PushToKafka("user_bookings", req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return respons
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
	fmt.Println(req.Response)
	AcquireDriverLock(req.DriverID)
	if req.Response == "ACCEPTED" {
		err := UpdateBookingStatus(req.BookingID, "ACCEPTED")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if req.Response == "REJECT" {
		ReleaseDriverLock(req.DriverID)
		err := UpdateBookingStatus(req.BookingID, "REJECTED")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if req.Response == "COMPLETED" {
		err := UpdateBookingStatus(req.BookingID, "COMPLETED")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		PushEventToRabbitMQ("booking_completed", req.BookingID)
		ReleaseDriverLock(req.DriverID)
	} else if req.Response == "STARTED" {
		fmt.Println("\n\nBooking started!!!!!!!\n\n")
		err := UpdateBookingStatus(req.BookingID, "STARTED")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		PushEventToRabbitMQ("booking_started", req.BookingID)
	}

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
		fmt.Println(err)
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
