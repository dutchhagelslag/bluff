package server

import (
	"bytes"
	"log"
	"time"
	"github.com/gorilla/websocket"
)

type Player struct {
	Room *Room              `json:"-"`
	Conn *websocket.Conn    `json:"-"`
	Send chan []byte        `json:"-"`
	Rdy bool                `json:"-"`

	Cards [2]Card
	Name string
	Coins int
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
		Room: nil,
		Conn: conn,
		Send: make(chan []byte),
		Name: name,
		Cards: [2]Card{},
		Coins: 2,
		Rdy: false,
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
	return (player.Cards[0] != Dead) && (player.Cards[1] != Dead)
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (player *Player) readPump() {
	defer func() {
		player.Room.remove(player.Name)
		player.Conn.Close()
	}()

	player.Conn.SetReadLimit(maxMessageSize)
	player.Conn.SetReadDeadline(time.Now().Add(pongWait))
	player.Conn.SetPongHandler(func(string) error { player.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := player.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		player.Room.Receive <- message
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
		player.Room.remove(player.Name)
		player.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-player.Send:
			player.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Players been kicked
				player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := player.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(player.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-player.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			player.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := player.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}


