package main

type BookingRequest struct {
	UserID          int       `json:"user_id"`
	PickupLocation  Location  `json:"pickup_location"`
	DropoffLocation Location  `json:"dropoff_location"`
	FareAmount      float64   `json:"fare_amount"`
}

type DriverResponse struct {
	DriverID   int    `json:"driver_id"`
	BookingID  int    `json:"booking_id"`
	Response   string `json:"response"` 
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
