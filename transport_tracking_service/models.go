package main

type Booking struct {
	ID              int
	UserID          int
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
