package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
)

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in KM
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := R * c
	return distance
}

func CheckSurgePricing(location Location) float64 {
	key := fmt.Sprintf("fare_requests:%f:%f", location.Lat, location.Lng)
	count, err := redisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return 1.0
	}
	if count > 5 {
		return 1.5 // 50% surge
	}
	return 1.0
}

func SaveFareRequestInRedis(userID int, location Location) {
	key := fmt.Sprintf("fare_requests:%f:%f", location.Lat, location.Lng)
	redisClient.Incr(ctx, key)
	redisClient.Expire(ctx, key, time.Minute*1)
}

func SavePayment(payment Payment) (int, error) {
	var paymentID int
	ctx := context.Background()
	err := dbPool.QueryRow(ctx, `
		INSERT INTO payments (user_id, pickup_location, dropoff_location, distance_km, fare_amount, status)
		VALUES ($1, ST_GeomFromText($2, 4326), ST_GeomFromText($3, 4326), $4, $5, $6) RETURNING id
	`, payment.UserID, payment.PickupLocation, payment.DropoffLocation, payment.DistanceKM, payment.FareAmount, payment.Status).Scan(&paymentID)
	return paymentID, err
}
