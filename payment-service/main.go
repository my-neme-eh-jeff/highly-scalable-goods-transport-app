package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/get-fare", GetFare).Methods("POST")

	handler := cors.Default().Handler(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, handler)

	log.Println("Payment Service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
