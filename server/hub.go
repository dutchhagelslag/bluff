package main

import(
	"strconv"
)

type Hub struct {
	off chan []byte
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
	room_id int
}

func newHub(room_id int) *Hub {
	return &Hub{
		off: make(chan []byte),
		room_id: room_id,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {

		// close hub and all players
		case <-h.off:
			for client := range h.clients {
				client.off <- []byte("off")
			}
			return

		case client := <-h.register:
			h.clients[client] = true
			client.send <- []byte(strconv.Itoa(h.room_id))

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
