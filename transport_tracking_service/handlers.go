// handlers.go

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
    activeBookings   = make(map[string]bool)
    activeBookingsMu sync.RWMutex
)

func addActiveBooking(bookingID string) {
    activeBookingsMu.Lock()
    defer activeBookingsMu.Unlock()
    activeBookings[bookingID] = true
}

func removeActiveBooking(bookingID string) {
    activeBookingsMu.Lock()
    defer activeBookingsMu.Unlock()
    delete(activeBookings, bookingID)
}

func isBookingActive(bookingID string) bool {
    activeBookingsMu.RLock()
    defer activeBookingsMu.RUnlock()
    return activeBookings[bookingID]
}

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

    // Check if booking is active
    if !isBookingActive(bookingID) {
        fmt.Fprintf(w, "event: end\ndata: Booking not active\n\n")
        flusher.Flush()
        return
    }

    // Send historical data
    locations, err := GetLocationsFromMongoDB(bookingID)
    if err != nil {
        http.Error(w, "Error fetching locations", http.StatusInternalServerError)
        return
    }
    for _, loc := range locations {
        locJSON, _ := json.Marshal(loc.Location)
        fmt.Fprintf(w, "data: %s\n\n", locJSON)
    }
    flusher.Flush()

    // Listen for new location updates
    changeStream, err := GetMongoChangeStream(bookingID)
    if err != nil {
        http.Error(w, "Error opening change stream", http.StatusInternalServerError)
        return
    }
    defer changeStream.Close(ctx)

    for {
        if !isBookingActive(bookingID) {
            fmt.Fprintf(w, "event: end\ndata: Booking completed\n\n")
            flusher.Flush()
            break
        }

        if changeStream.Next(ctx) {
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
        } else if err := changeStream.Err(); err != nil {
            log.Printf("Change stream error: %v", err)
            break
        }
    }
}
