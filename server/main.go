package main

import (
    "strconv"
	"fmt"
	"bluff/server/spinner"
    "github.com/julienschmidt/httprouter"
    "log"
    "net/http"
	"sync"
)

var all_rooms sync.Map

func main(){
	fmt.Println("hi")
    health_check := make(chan string)

    go spinner.Spin(health_check)

    router := httprouter.New()
    router.GET("/create_room/:player_name", create_room)
    router.GET("/join_room/:room_id/:player_name", join_room)

    router.GET("/draw/:room_id", test)

    router.GET("/rm/:room_id/:player_name", rm)


    log.Fatal(http.ListenAndServe(":8080", router))
}

func rm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("room_id"))
	p := ps.ByName("player_name")

	room_void, ok := all_rooms.Load(id)
	if !ok {
		fmt.Println("Failed to find room")
		return
	}
	room_ptr := room_void.(*room)

	fmt.Println(room_ptr.id)
	fmt.Println(p)
	remove_player(p, room_ptr)

	print_rooms("rm")
	return
}

func test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("room_id"))
	room_void, _ := all_rooms.Load(id)
	room_ptr := room_void.(*room)

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

	room_id := global_id_get()
	host_name := ps.ByName("player_name")

	new_room := init_room(host_name, room_id)
	all_rooms.Store(room_id, new_room)

	init_client(conn, new_room.hub)

	print_rooms("create room")
	return
}

func join_room(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	player_name := ps.ByName("player_name")

	id, _ := strconv.Atoi(ps.ByName("room_id"))

	room_void, ok := all_rooms.Load(id)

	// room doesn't exist
	if(!ok){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	room_ptr := room_void.(*room)

	// room is full
	if(len(room_ptr.members) == cap(room_ptr.members)){
		w.WriteHeader(http.StatusConflict)
		return
	}

	// game has already started
	if(room_ptr.turn != 0){
		w.WriteHeader(http.StatusConflict)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	new_player := player{
		name: player_name,
		cards: [2]card{},
		coins: 2,
	}

	room_ptr.members = append(room_ptr.members,new_player)

	w.WriteHeader(http.StatusOK)

	init_client(conn, room_ptr.hub)

	print_rooms("join_room")
}

func start_game(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	id, _ := strconv.Atoi(ps.ByName("room_id"))

	room_void, ok := all_rooms.Load(id)
	room_ptr := room_void.(*room)

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

