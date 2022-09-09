package main

import(
	"strconv"
	"bytes"
)

type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
	room_id int
}

func newHub(room_id int) *Hub {
	return &Hub{
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
				if bytes.Compare(message, []byte("kill-all")) == 0 {
					close(client.send)
					delete(h.clients, client)
				}else{
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}

			if bytes.Compare(message, []byte("kill-all")) == 0 {
				return
			}
		}
	}
}
