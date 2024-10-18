package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
    hub := NewHub()
    go hub.Run()

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    consumer, err := NewKafkaConsumer(hub)
    if err != nil {
        log.Fatalf("Failed to create Kafka consumer: %v", err)
    }

    // Start consumer with reconnection logic
    go func() {
        for {
            consumer.Start(ctx) 
            select {
            case <-ctx.Done():
                return
            default:
                log.Printf("Kafka consumer restarting in 5 seconds...")
                time.Sleep(time.Second * 10)
            }
        }
    }()

    router := mux.NewRouter()

    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        if consumer.IsHealthy() {
            w.WriteHeader(http.StatusOK)
            json.NewEncoder(w).Encode(consumer.GetMetrics())
        } else {
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "status": "unhealthy",
                "metrics": consumer.GetMetrics(),
            })
        }
    })

    router.HandleFunc("/ws/driver/assign", func(w http.ResponseWriter, r *http.Request) {
        ServeDriverWS(hub, w, r)
    })

    handler := cors.Default().Handler(router)

    server := &http.Server{
        Addr:    ":8084",
        Handler: handler,
    }

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    go func() {
        log.Println("Transport Matcher Service running on port 8084")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error starting server: %v", err)
        }
    }()

    // Wait for interruption signal
    <-stop
    log.Println("Shutting down server...")

    // Graceful shutdown
    cancel()
    if err := server.Shutdown(context.Background()); err != nil {
        log.Printf("Error shutting down server: %v", err)
    }
}

