package server

import(
	"log"
	"testing"
	"net/url"
	"github.com/gorilla/websocket"
)

func TestLobby(t *testing.T){
	log.SetFlags(0)

	// setup server
	go RunServer()


	t.Run("create new lobby", func(t *testing.T) {
		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/create_room/tester"}

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal(err)
			t.Errorf("Failed creating websocket client")
		}
		defer c.Close()
	})





	// make calls and test values

	// if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
	// 	t.Errorf("Failed to write over websocket")
	// }

	// var msg = make([]byte, 512)

	// var n int
	// if n, err = ws.Read(msg); err != nil {
	// 	log.Fatal(err)
	// }
}
