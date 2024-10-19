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
    defer Cleanup() 

    router := mux.NewRouter()

    router.HandleFunc("/api/user/book-transport", BookTransport).Methods("POST")
    router.HandleFunc("/api/driver/respond-booking", DriverRespondBooking).Methods("POST")
    router.HandleFunc("/api/user/bookings", GetUserBookings).Methods("GET")
    router.HandleFunc("/api/driver/complete-ride", DriverCompleteRide).Methods("POST")

    handler := cors.Default().Handler(router)
    loggedRouter := handlers.LoggingHandler(os.Stdout, handler)

    log.Println("Transport Service running on port 8081")
    log.Fatal(http.ListenAndServe(":8081", loggedRouter))
}
