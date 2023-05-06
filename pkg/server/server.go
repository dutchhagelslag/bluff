package server

import (
	"log"
	"fmt"
	"sync"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
)

var all_rooms sync.Map

type Page struct {
	Title string
	Body []byte
}

func run_server(){
	router := httprouter.New()

	router.GET("/mping", mping)

	// router.GET("/create_room/:player_name", create_room)
	// router.GET("/join_room/:room_id/:player_name", join_room)
	// router.GET("/draw/:room_id", test)
	// router.GET("/rm/:room_id/:player_name", rm)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func mping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "learning something i suppose")
}

func rm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("room_id"))
	player := ps.ByName("player_name")

	room_void, ok := all_rooms.Load(id)

	if !ok {
		fmt.Println("Failed to find room")
		return
	}

	room := room_void.(*Room)

	fmt.Println(room.Id)
	fmt.Println(player)

	room.remove(player)

	print_rooms("rm")
	return
}

func test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("room_id"))
	room_void, _ := all_rooms.Load(id)
	room_ptr := room_void.(*Room)

	draw_cards(room_ptr)

	print_rooms("test")
	return
}

func create_room(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// setup player connection //
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	room_id, err := uuid.NewV4()

	if err != nil {
		return // error generating uuid
	}
	host_name := ps.ByName("player_name")

	new_room := init_room(room_id.String())

	new_player := init_player(host_name, conn)

	new_room.add(new_player)
	print_rooms("create room")
	return
}

func join_room(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("room_id"))

	room_void, ok := all_rooms.Load(id)

	if !ok{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	room := room_void.(*Room)


	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	player_name := ps.ByName("player_name")

	if string(player_name)[0] == '!'{
		fmt.Println("Invalid player name: cant start with !")
	}

	new_player := init_player(player_name, conn)

	room.add(new_player)
	// if err := room.add_member(new_player); err != nil{
	// 	w.WriteHeader(http.StatusConflict)
	// 	return
	// }

	print_rooms("join_room")
}

func start_game(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id, _ := strconv.Atoi(ps.ByName("room_id"))

	room_void, ok := all_rooms.Load(id)
	room_ptr := room_void.(*Room)

	// todo: verify host
	// host_name := ps.ByName("host_name")

	// room doesn't exist
	if(!ok){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// start game state
	init_game(room_ptr)


	// todo setup websockets with players
	print_rooms("start game")
}


func init_game(lobby *Room){
	// change mode of hub
	draw_cards(lobby)
	lobby.Turn += 1
}

// all important debug print
func print_rooms(action string){
	fmt.Println("======================================================================")
	fmt.Println(action)

	m := map[int]interface{}{}

	all_rooms.Range(func(key, value interface{}) bool {
		body := value.(*Room)
		room_output := ""
		hands := " || hands: "

		room_output += "Host: " + body.Members[0].Name + " || Members:"
		for i := range body.Members{
			room_output += " " + body.Members[i].Name
		}
		for i := range body.Members{
			hands += " " + body.Members[i].Name + ": "
			hands += body.Members[i].Cards[0].String() + " and "
			hands += body.Members[i].Cards[1].String() + "||"
		}

		// room_output["players"] = body.owner.Name + room_mems

		// a, _ := json.Marshal(room_output)

		deck, _ := json.MarshalIndent(body.Deck, "", " ")
		fmt.Println(string(deck))

		m[key.(int)] = string(room_output + hands)
		fmt.Println("LobbyLen: " + strconv.Itoa(len(body.Members)))
		fmt.Println("LobbyCap: " + strconv.Itoa(cap(body.Members)))

		return true
	})

	b, _ := json.MarshalIndent(m, "", " ")
	fmt.Println(string(b))
	fmt.Println("======================================================================")
}


// checklist

// func turn_check(){

// }



// Card nulls ->

// watch state

// process turn
  // check if player is alive -> skip turn / display
  // check if game is over -> is more than one people alive -> trigger game end sequence
  // Check if action is valid -> is player bluffing
  // turn timer

  //


  // while true
  //   process-turn-> turn % player count // test ing

  //

  // back to lobby state
