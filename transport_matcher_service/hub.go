package main

import (
	"log"
)

type Hub struct {
	// Registered clients.
	clients map[int]*Client 
	
	// Inbound messages from the clients.
	broadcast chan interface{}

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]*Client),
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.driverID] = client
			log.Printf("Driver %d connected.", client.driverID)
		case client := <-h.unregister:
			if _, ok := h.clients[client.driverID]; ok {
				delete(h.clients, client.driverID)
				close(client.send)
				log.Printf("Driver %d disconnected.", client.driverID)
			}
		// case message := <-h.broadcast:

		}
	}
}

func (h *Hub) SendToDriver(driverID int, message interface{}) {
	if client, ok := h.clients[driverID]; ok {
		client.send <- message
	} else {
		log.Printf("Driver %d is not connected.", driverID)
	}
}
