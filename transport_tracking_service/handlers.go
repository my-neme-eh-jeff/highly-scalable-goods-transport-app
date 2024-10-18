package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func TrackTransport(w http.ResponseWriter, r *http.Request) {
	bookingID := r.URL.Query().Get("booking_id")
	if bookingID == "" {
		http.Error(w, "Missing booking_id parameter", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Send initial data
	locations := GetLocationsFromMongoDB(bookingID)
	for _, loc := range locations {
		locJSON, _ := json.Marshal(loc)
		fmt.Fprintf(w, "data: %s\n\n", locJSON)
	}
	flusher.Flush()

	// Listen for new location updates
	changeStream := GetMongoChangeStream(bookingID)

	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var event struct {
			FullDocument LocationRecord `bson:"fullDocument"`
		}
		if err := changeStream.Decode(&event); err != nil {
			log.Printf("Error decoding change stream document: %v", err)
			continue
		}

		locJSON, _ := json.Marshal(event.FullDocument.Location)
		fmt.Fprintf(w, "data: %s\n\n", string(locJSON))
		flusher.Flush()
	}

	if err := changeStream.Err(); err != nil {
		log.Printf("Change stream error: %v", err)
	}
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

	json.NewEncoder(w).Encode(bookings)
}
func GetBookingsByUserID(userID int64) ([]BookingDetailsResponse, error) {
	var bookings []BookingDetailsResponse

	query := `
        SELECT id, user_id, driver_id,
               ST_Y(pickup_location::geometry) AS pickup_lat,
               ST_X(pickup_location::geometry) AS pickup_lng,
               ST_Y(dropoff_location::geometry) AS dropoff_lat,
               ST_X(dropoff_location::geometry) AS dropoff_lng,
               fare_amount, status
        FROM bookings
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := dbPool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking BookingDetailsResponse
		err := rows.Scan(
			&booking.BookingID,
			&booking.UserID,
			&booking.DriverID,
			&booking.PickupLocation.Lat,
			&booking.PickupLocation.Lng,
			&booking.DropoffLocation.Lat,
			&booking.DropoffLocation.Lng,
			&booking.FareAmount,
			&booking.Status,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
