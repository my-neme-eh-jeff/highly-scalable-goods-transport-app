package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/api/user/track-transport", TrackTransport).Methods("GET")

    // Start a goroutine to listen to RabbitMQ events
    go ListenToRabbitMQ()

    handler := cors.Default().Handler(router)
    loggedRouter := handlers.LoggingHandler(os.Stdout, handler)

    log.Println("Transport Tracking Service running on port 8082")
    log.Fatal(http.ListenAndServe(":8082", loggedRouter))
}
