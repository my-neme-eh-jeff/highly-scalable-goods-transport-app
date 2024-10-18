package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	hub      *Hub
	topics   []string
	metrics  KafkaMetrics
}

type KafkaMetrics struct {
	MessagesProcessed uint64
	ErrorCount        uint64
	LastError         string
	LastErrorTime     time.Time
	Connected         bool
}

func (k *KafkaConsumer) Start(ctx context.Context) {
	log.Printf("Starting Kafka consumer for topics: %v", k.topics)

	err := k.consumer.SubscribeTopics(k.topics, nil)
	if err != nil {
		log.Printf("Failed to subscribe to topics: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer stopping...")
			k.consumer.Close()
			return

		default:
			msg, err := k.consumer.ReadMessage(time.Second * 1)

			if err != nil {
				kafkaErr, isKafkaErr := err.(kafka.Error)

				// Handle different types of errors
				if isKafkaErr {
					// Timeout is normal, just continue
					if kafkaErr.Code() == kafka.ErrTimedOut {
						continue
					}

					// Connection errors
					if kafkaErr.Code() == kafka.ErrTransport ||
						kafkaErr.Code() == kafka.ErrBrokerNotAvailable {
						k.metrics.Connected = false
						k.recordError(err)
						log.Printf("Connection error: %v. Attempting to reconnect...", err)
						time.Sleep(time.Second * 5)
						continue
					}
				}

				// Record other errors
				k.recordError(err)
				log.Printf("Error reading message: %v", err)
				continue
			}

			// Mark as connected on successful message read
			k.metrics.Connected = true

			// Process message
			if err := k.processMessage(msg); err != nil {
				k.recordError(err)
				log.Printf("Error processing message: %v", err)
				continue
			}

			atomic.AddUint64(&k.metrics.MessagesProcessed, 1)
		}
	}
}

func (k *KafkaConsumer) tryProcessBooking(booking BookingRequest) error {
	// Find nearby drivers
	drivers, err := FindNearbyDrivers(booking.PickupLocation)
	if err != nil {
		return fmt.Errorf("failed to find nearby drivers: %v", err)
	}

	bookingID := booking.BookingID // Use the booking ID from the message

	// Try to assign a driver
	for _, driver := range drivers {
		if AcquireDriverLock(driver.DriverID) {
			err := AssignDriver(bookingID, driver.DriverID)
			if err != nil {
				ReleaseDriverLock(driver.DriverID)
				continue
			}

			// Create the complete booking object for notification
			completeBooking := Booking{
				BookingID:       bookingID,
				UserID:          booking.UserID,
				PickupLocation:  booking.PickupLocation,
				DropoffLocation: booking.DropoffLocation,
				FareAmount:      booking.FareAmount,
				Status:          "DRIVER_ASSIGNED",
			}

			// Notify the driver
			NotifyDriver(k.hub, driver.DriverID, completeBooking)
			return nil
		}
	}

	return fmt.Errorf("no available drivers found")
}

func (k *KafkaConsumer) processMessage(msg *kafka.Message) error {
	var booking BookingRequest
	if err := json.Unmarshal(msg.Value, &booking); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}

	// Process with retries
	return retry(3, time.Second, func() error {
		return k.tryProcessBooking(booking)
	})
}

func NewKafkaConsumer(hub *Hub) (*KafkaConsumer, error) {
	config := &kafka.ConfigMap{
		// Connection configs from your client.properties
		"bootstrap.servers": "pkc-7prvp.centralindia.azure.confluent.cloud:9092",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     "YBFJ64E5XT4BCQUN",
		"sasl.password":     "4rG5gE/QfHk+WwiCbTKSgKXoJ4ZJytu7kGKaqej/pm77CVfdMfXwt7rKeng4sied",

		// Group configuration
		"group.id":           "transport-matcher-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,

		// Availability and resilience settings
		"session.timeout.ms":    45000,  // From your properties file
		"heartbeat.interval.ms": 15000,  // One-third of session timeout
		"max.poll.interval.ms":  300000, // 5 minutes

		// Reconnection settings
		"socket.keepalive.enable":  true,
		"reconnect.backoff.ms":     1000,  // Initial backoff
		"reconnect.backoff.max.ms": 30000, // Max backoff

		// Request handling
		"retry.backoff.ms":   1000,
		"request.timeout.ms": 30000,

		// Performance optimization
		"fetch.min.bytes": 1,

		// Debug settings (optional - remove in production)
		// "debug": "consumer,cgrp,topic,fetch",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		hub:      hub,
		topics:   []string{"user_bookings"}, // adjust your topic name
	}, nil
}

func (k *KafkaConsumer) recordError(err error) {
	atomic.AddUint64(&k.metrics.ErrorCount, 1)
	k.metrics.LastError = err.Error()
	k.metrics.LastErrorTime = time.Now()
}

// Helper retry function
func retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, sleep*2, fn)
		}
		return err
	}
	return nil
}

// Health check method
func (k *KafkaConsumer) IsHealthy() bool {
	if !k.metrics.Connected {
		return false
	}

	// Check if we've had any errors in the last minute
	if k.metrics.LastErrorTime.After(time.Now().Add(-time.Minute)) {
		return false
	}

	return true
}

// Get metrics
func (k *KafkaConsumer) GetMetrics() KafkaMetrics {
	return k.metrics
}
