package main

import (
    "strconv"
	"fmt"
    "github.com/julienschmidt/httprouter"
    "log"
    "net/http"
	"sync"
	"html/template"
	"os"
)

var all_rooms sync.Map

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func main(){
    health_check := make(chan string)

    go spin(health_check)

    router := httprouter.New()

    router.GET("/hi", testw)
    // router.GET("/create_room/:player_name", create_room)
    // router.GET("/join_room/:room_id/:player_name", join_room)
    // router.GET("/draw/:room_id", test)
    // router.GET("/rm/:room_id/:player_name", rm)

    log.Fatal(http.ListenAndServe(":8080", router))
}

func testw(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := &Page{Title: "test", Body: []byte("hi")}
	t, _ := template.ParseFiles("view.html")
	t.Execute(w,p)
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

	room_id := global_id_get()
	host_name := ps.ByName("player_name")

	new_room := init_room(room_id)
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

