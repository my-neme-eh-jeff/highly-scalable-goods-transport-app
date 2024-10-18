package main

import (
    "log"
    "os"

    "github.com/streadway/amqp"
)

func ListenToRabbitMQ() {
    conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    bookingAcceptedQueue, err := ch.QueueDeclare(
        "booking_accepted", // name
        true,               // durable
        false,              // delete when unused
        false,              // exclusive
        false,              // no-wait
        nil,                // arguments
    )
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    bookingCompletedQueue, err := ch.QueueDeclare(
        "booking_completed",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    // Consume messages from the queues
    acceptedMsgs, err := ch.Consume(
        bookingAcceptedQueue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to register consumer: %v", err)
    }

    completedMsgs, err := ch.Consume(
        bookingCompletedQueue.Name,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Failed to register consumer: %v", err)
    }

    log.Println("Listening to RabbitMQ...")
    // Start goroutines to handle messages
    go handleAcceptedBookings(acceptedMsgs)
    go handleCompletedBookings(completedMsgs)

    select {}
}

func handleAcceptedBookings(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        bookingID := string(d.Body)
        log.Printf("Received booking accepted event for booking ID: %s", bookingID)
        addActiveBooking(bookingID)
    }
}

func handleCompletedBookings(msgs <-chan amqp.Delivery) {
    for d := range msgs {
        bookingID := string(d.Body)
        log.Printf("Received booking completed event for booking ID: %s", bookingID)
        removeActiveBooking(bookingID)
    }
}
