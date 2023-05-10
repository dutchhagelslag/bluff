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

	lobbies := []struct {
		host string
		members []string
		roomID string
	}{
		{ host: "ryan", members: []string{"hello", "xxx_green_xxx", "edgyboy101", "patrick"} },

		{ host: "jiyoung", members: []string{"sadboi202", "timetime", "patrick"} },

		{ host: "evan", members: []string{"micheal", "maggie", "jacob", "patrick"} },

		{ host: "rifah", members: []string{"frech", "bengal", "hui", "patrick"} },
	}

	// Create Lobbies and Add members
	for i, lobby := range lobbies {
		var postLobby sync.WaitGroup
		var postCreate sync.WaitGroup

		member_count := len(lobby.members)
		postLobby.Add(member_count)
		postCreate.Add(1)

		go t.Run(fmt.Sprintf("Create Lobby %d: %s", i, lobby.host), func(t *testing.T) {

			// or maby dont close all for more testing?
			defer fmt.Printf("Closing Lobby %d: %s", i, lobby.host)

			path := fmt.Sprintf("/create_room/%s", lobby.host)
			fmt.Println(path)

			u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: path}

			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

			if err != nil {
				log.Fatal(err)
				t.Errorf("Failed creating websocket client")
			}
			defer conn.Close() // same as above

			// Receive a message from the WebSocket server
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
				t.Errorf("Error receiving message from server:")
			}
			log.Println("Room ID:", string(message))
			lobby.roomID = string(message)

			postCreate.Done()
			postLobby.Wait()
		})


		for _, member := range lobby.members {
			t.Run(fmt.Sprintf("Add Members %d: %s", i, member), func(t *testing.T) {
				defer postLobby.Done()

				postCreate.Wait()

				path := fmt.Sprintf("/join_room/%s/%s", lobby.roomID, member)
				fmt.Println(path)

				u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: path}

				_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					t.Errorf("Failed creating websocket client")
				}

				// defer conn.Close()
			})
		}
	}


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
