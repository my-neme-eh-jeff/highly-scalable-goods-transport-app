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

	router.HandleFunc("/api/match-drivers", MatchDrivers).Methods("POST")

	handler := cors.Default().Handler(router)
	loggedRouter := handlers.LoggingHandler(os.Stdout, handler)

	log.Println("Transport Matcher Service running on port 8084")
	log.Fatal(http.ListenAndServe(":8084", loggedRouter))
}
