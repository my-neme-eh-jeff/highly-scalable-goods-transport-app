	package main

	import (
		"log"
		"time"

		"github.com/gorilla/websocket"
	)

	const (
		writeWait  = 10 * time.Second
		pongWait   = 60 * time.Second
		pingPeriod = (pongWait * 9) / 10
	)

	type Client struct {
		driverID int
		hub      *Hub
		conn     *websocket.Conn
		send     chan interface{}
	}

	func (c *Client) readPump() {
		defer func() {
			c.hub.unregister <- c
			c.conn.Close()
		}()
		c.conn.SetReadLimit(512)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.conn.SetPongHandler(func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
		for {
			_, _, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}
			// Process incoming messages if needed
		}
	}

	func (c *Client) writePump() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			c.conn.Close()
		}()
		for {
			select {
			case message, ok := <-c.send:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					c.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}
				if err := c.conn.WriteJSON(message); err != nil {
					log.Printf("WriteJSON error: %v", err)
					return
				}
			case <-ticker.C:
				c.conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
