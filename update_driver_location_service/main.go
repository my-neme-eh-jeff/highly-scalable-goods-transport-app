package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {

    router := mux.NewRouter()

    router.HandleFunc("/ws/driver/update-location", DriverLocationWebSocket).Methods("GET")

    log.Println("Update Driver Location Service running on port 8083")
    log.Fatal(http.ListenAndServe(":8083", router))
}
