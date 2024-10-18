package main

type Booking struct {
	BookingID       int      `json:"booking_id"`
	UserID          int64      `json:"user_id"`
	PickupLocation  Location `json:"pickup_location"`
	DropoffLocation Location `json:"dropoff_location"`
	FareAmount      float64  `json:"fare_amount"`
	Status          string   `json:"status"`
}

type Driver struct {
	DriverID int     `json:"driver_id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type BookingRequest struct {
    BookingID       int       `json:"booking_id"`
    UserID          int64            `json:"user_id"`
    PickupLocation  Location  `json:"pickup_location"`
    DropoffLocation Location  `json:"dropoff_location"`
    FareAmount      float64   `json:"fare_amount"`
}

