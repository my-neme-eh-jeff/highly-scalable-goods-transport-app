package main

import "time"

type LocationRecord struct {
	BookingID string    `bson:"booking_id"`
	Location  Location  `bson:"location"`
	Timestamp time.Time `bson:"timestamp"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

