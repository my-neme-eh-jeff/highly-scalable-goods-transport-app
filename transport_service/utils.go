package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	dbPool      *pgxpool.Pool
)

func init() {
	// Load environment variables
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Printf("Error loading .env.local file: %v", err)
	}

	// Redis initialization
	redisOptions, err := redis.ParseURL(os.Getenv("REDIS_UPSTASH_ADDR"))
	if err != nil {
		log.Fatalf("Error parsing Redis URL: %v", err)
	}
	redisClient = redis.NewClient(redisOptions)

	// Postgres connection pool
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("NEON_DB_URL"))
	if err != nil {
		log.Fatalf("Unable to parse Postgres config: %v", err)
	}
	dbPool, err = pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Test connections
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Redis connection error: %v", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Fatalf("Postgres connection error: %v", err)
	}

	// Ensure the 'bookings' table exists
	_, err = dbPool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS bookings (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            driver_id INTEGER,
            pickup_location GEOGRAPHY(POINT),
            dropoff_location GEOGRAPHY(POINT),
            fare_amount DECIMAL(10,2),
            status VARCHAR(20),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatalf("Error creating 'bookings' table: %v", err)
	}

	log.Println("Connected to Redis and Postgres, and ensured 'bookings' table exists.")
}

func SaveBooking(req BookingRequest) (int, error) {
	var bookingID int
	query := `
        INSERT INTO bookings (user_id, pickup_location, dropoff_location, fare_amount, status)
        VALUES ($1, ST_GeographyFromText($2), ST_GeographyFromText($3), $4, $5) RETURNING id`
	err := dbPool.QueryRow(ctx, query,
		req.UserID,
		fmt.Sprintf("SRID=4326;POINT(%f %f)", req.PickupLocation.Lng, req.PickupLocation.Lat),
		fmt.Sprintf("SRID=4326;POINT(%f %f)", req.DropoffLocation.Lng, req.DropoffLocation.Lat),
		req.FareAmount,
		"REQUESTED",
	).Scan(&bookingID)

	if err != nil {
		log.Printf("Error saving booking: %v", err)
		return 0, err
	}
	return bookingID, nil
}

func PushToKafka(topic string, message interface{}) error {
	config := ReadKafkaConfig()
	log.Printf("Kafka Config: %v", config)
	p, err := kafka.NewProducer(&config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %v", err)
	}
	defer p.Close()

	// Goroutine for handling message delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced to topic %s: key = %-10s value = %s\n", *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	msgBytes, _ := json.Marshal(message)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Timestamp:      time.Now(),
		Value:          msgBytes,
	}, nil)

	
	// Flush messages to Kafka
	p.Flush(15000)
	return err
}

func UpdateBookingStatus(bookingID int, status string) error {
    _, err := dbPool.Exec(ctx, `
        UPDATE bookings SET status = $1 WHERE id = $2
    `, status, bookingID)
    if err != nil {
        log.Printf("Error updating booking status: %v", err)
        return err
    }
    return nil
}


func PushEventToRabbitMQ(eventType string, bookingID int) error {
	conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(eventType, false, true, false, false, nil)
	if err != nil {
		return err
	}

	body := fmt.Sprintf("%d", bookingID)
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		return err
	}
	log.Printf("Pushed event %s with bookingID %d to RabbitMQ", eventType, bookingID)
	return nil
}

func AcquireDriverLock(driverID int) bool {
	lockKey := fmt.Sprintf("driver_lock:%d", driverID)
	success, err := redisClient.SetNX(ctx, lockKey, "locked", time.Minute).Result()
	if err != nil {
		log.Printf("Failed to acquire driver lock: %v", err)
		return false
	}
	return success
}

func ReleaseDriverLock(driverID int) {
	lockKey := fmt.Sprintf("driver_lock:%d", driverID)
	err := redisClient.Del(ctx, lockKey).Err()
	if err != nil {
		log.Printf("Failed to release driver lock: %v", err)
	}
}

func GetBookingsByUserID(userID int64) ([]Booking, error) {
	var bookings []Booking
	query := `
        SELECT id, user_id, driver_id,
               ST_Y(pickup_location::geometry) AS pickup_lat,
               ST_X(pickup_location::geometry) AS pickup_lng,
               ST_Y(dropoff_location::geometry) AS dropoff_lat,
               ST_X(dropoff_location::geometry) AS dropoff_lng,
               fare_amount, status, created_at
        FROM bookings
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
	rows, err := dbPool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var booking Booking
		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.DriverID,
			&booking.PickupLocation.Lat,
			&booking.PickupLocation.Lng,
			&booking.DropoffLocation.Lat,
			&booking.DropoffLocation.Lng,
			&booking.FareAmount,
			&booking.Status,
			&booking.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
