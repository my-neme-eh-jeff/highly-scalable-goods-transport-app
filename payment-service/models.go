package main

import "time"

type FareRequest struct {
	UserID          int      `json:"user_id"`
	PickupLocation  Location `json:"pickup_location"`
	DropoffLocation Location `json:"dropoff_location"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Payment struct {
	ID              int
	UserID          int
	PickupLocation  string
	DropoffLocation string
	DistanceKM      float64
	FareAmount      float64
	Status          string
	CreatedAt       time.Time
}
