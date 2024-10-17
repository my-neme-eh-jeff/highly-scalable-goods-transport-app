package main

import (
	"fmt"
	"net/http"
	"time"
)

func TrackTransport(w http.ResponseWriter, r *http.Request) {
	bookingID := r.URL.Query().Get("booking_id")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Send historical data
	locations := GetLocationsFromMongoDB(bookingID)
	for _, loc := range locations {
		fmt.Fprintf(w, "data: %s\n\n", loc)
	}
	flusher.Flush()

	// Send live updates
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			newLoc := GetLatestLocation(bookingID)
			fmt.Fprintf(w, "data: %s\n\n", newLoc)
			flusher.Flush()
		}
	}
}
