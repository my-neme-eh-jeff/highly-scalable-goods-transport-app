package main

import "time"

type Booking struct {
	ID              int
	UserID          int64
	DriverID        int
	PickupLocation  Location
	DropoffLocation Location
	FareAmount      float64
	Status          string
	CreatedAt       string
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type BookingDetailsResponse struct {
	BookingID       int      `json:"booking_id"`
	UserID          int64    `json:"user_id"`
	DriverID        int      `json:"driver_id"`
	PickupLocation  Location `json:"pickup_location"`
	DropoffLocation Location `json:"dropoff_location"`
	FareAmount      float64  `json:"fare_amount"`
	Status          string   `json:"status"`
}

type LocationRecord struct {
	BookingID string    `bson:"booking_id"`
	Location  Location  `bson:"location"`
	Timestamp time.Time `bson:"timestamp"`
}

