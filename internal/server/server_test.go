package server

import(
	"fmt"
	"log"
	"sync"
	"testing"
	"net/url"
	"github.com/gorilla/websocket"
)

func TestLobby(t *testing.T){
	log.SetFlags(0)

	// setup server
	go RunServer()

	// var roomID string

	var wg sync.WaitGroup
	wg.Add(1)



	t.Run("create new lobby", func(t *testing.T) {
		defer fmt.Println("created new lobby")

		path := "/create_room/tester"
		fmt.Println(path)

		u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: path}

		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal(err)
			t.Errorf("Failed creating websocket client")
		}
		defer conn.Close()

		// Receive a message from the WebSocket server
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			t.Errorf("Error receiving message from server:")
		}
		log.Println("Room ID:", string(message))
		// roomID = string(message)

	})




	t.Run("join lobby", func(t *testing.T) {
		wg.Wait()
		// path := fmt.Sprintf("/join_room/%s/patrick", roomID)
		// fmt.Println(path)

		fmt.Println("hihih")

		// u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: path}

		// conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		// if err != nil {
		// 	t.Errorf("Failed creating websocket client")
		// }

		// defer conn.Close()
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
