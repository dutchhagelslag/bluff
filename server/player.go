package main

import (
	"bytes"
	"log"
	"time"
	"github.com/gorilla/websocket"

)

type Player struct {
	room *Room              `json:"-"`
	conn *websocket.Conn    `json:"-"`
	send chan []byte        `json:"-"`

    cards [2]Card
	rdy bool
    name string
    coins int
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func init_player(name string, conn *websocket.Conn) *Player{
	new_player := &Player{
		room: nil,
		conn: conn,
		send: make(chan []byte),
		name: name,
		cards: [2]Card{},
		coins: 2,
		rdy: false,
	}
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go new_player.writePump()
	go new_player.readPump()

	return new_player
}


func (player *Player) lose_card(lost_card Card){

}

func (player *Player) is_alive() bool{
	return (player.cards[0] != Dead) && (player.cards[1] != Dead)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (player *Player) readPump() {
	defer func() {
		player.room.remove(player.name)
		player.conn.Close()
	}()

	player.conn.SetReadLimit(maxMessageSize)
	player.conn.SetReadDeadline(time.Now().Add(pongWait))
	player.conn.SetPongHandler(func(string) error { player.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := player.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		player.room.receive <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (player *Player) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		player.room.remove(player.name)
		player.conn.Close()
	}()
	for {
		select {
		case message, ok := <-player.send:
			player.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Players been kicked
				player.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := player.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(player.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-player.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			player.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := player.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}


