package main

import "time"

type Location struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

type LocationRecord struct {
	BookingID string    `json:"booking_id" bson:"booking_id"`
	Location  Location  `json:"location" bson:"location"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
