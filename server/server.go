package main

import (
    "strconv"
    "fmt"
    // "os"
    "net/http"
    "log"
    "time"
    "regexp"
    "errors"
	"sync"
	"bluff/server/spinner"
	"encoding/json"
    "github.com/julienschmidt/httprouter"
    "math/rand"
)

// assume non malicous players for now
type id    string
type card int8

const(
    Dead card = iota
    Duke
    Assasin
    Contessa
    Captain
    Ambassador
)

type room struct{
    members []player
    turn int
    deck map[card]int
}

type player struct{
    name string
    cards [2]card
    coins int
	// member_of *room
	// websocket connection -> when disconnected remove()
}

func remove_player(remove_name string, lobby *room){
	for i := range lobby.members {
		player := lobby.members[i]
		fmt.Println(player.name)

		if(player.name == remove_name){
			// shift everyone
			lobby.members = append(lobby.members[:i], lobby.members[i+1:] ...)
			print_rooms()
			return
		}
	}
}

// randomly determine cards for players
func draw_cards(lobby *room){
	for i := 0 ; i < len(lobby.members) ; i++ {
		lobby.members[i].cards[0] = random_card(lobby.deck)
		lobby.members[i].cards[1] = random_card(lobby.deck)
	}
}

func random_card(deck map[card]int) card{
	var new_card card;
	for {

		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		new_card = card(random.Intn(5) + 1)

		if val, ok := deck[new_card]; (ok && val > 0){
			return new_card
		}
	}
}

// add validation

// randomly get new owner when owner leaves

// return current game_state: lobby_id

// return current lobby: lobby_id

// needs improvement

var validPath = regexp.MustCompile("^/(home|lobby|game)/([a-zA-Z0-9]+)$")

func validate_url(path string) (string, error){
    valid_path := validPath.FindStringSubmatch(path)
    if (valid_path == nil){
        return "", errors.New("invalid path")
    }
    return valid_path[2], nil
}

// 3 of each cards - cards in players hands (ignore till ambassador impl)
func init_deck() map[card]int{
	return map[card]int{
		Duke:3,
		Assasin:3,
		Contessa:3,
		Captain:3,
		Ambassador:3,
	}
}


func init_room(host_name string) *room{
	members_init := make([]player, 1, 6)

	// init host
	members_init[0] = player{
		name: host_name,
		cards: [2]card{},
		coins: 2,
	}

	return &room{
		members: members_init,
		turn: 0,
		deck: init_deck(),
	}
}

func init_game(lobby *room){
	lobby.turn += 1

}


var all_rooms sync.Map

func main(){
	fmt.Println("hi")
    health_check := make(chan string)

    go spinner.Spin(health_check)

    router := httprouter.New()
    router.GET("/create_room/:player_name", create_room)
    router.GET("/join_room/:room_id/:player_name", join_room)


    log.Fatal(http.ListenAndServe(":8080", router))
}

// until this games omega popular and needs an ID generator server
var global_id_counter = 0

func global_id_get() int {
	global_id_counter++
	return global_id_counter-1
}

// func turn_check(){

// }

func get_hand(){

}


func create_room(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)

	room_id := global_id_get()
	host_name := ps.ByName("player_name")

	// new room make, assume unique
	all_rooms.Store(room_id, init_room(host_name))

	// package room id for sending
	resp := make(map[string]string)
	resp["room_id"] = strconv.Itoa(room_id)

	json_resp, err := json.Marshal(resp)

	if(err != nil){
		log.Fatalf("Failed JSON marshal")
	}

	w.Write(json_resp)
	print_rooms()
	return
}

func join_room(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	player_name := ps.ByName("player_name")

	id, _ := strconv.Atoi(ps.ByName("room_id"))
	room_void, ok := all_rooms.Load(id)
	room_ptr := room_void.(*room)

	// room doesn't exist
	if(!ok){
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

	new_player := player{
		name: player_name,
		cards: [2]card{},
		coins: 2,
	}

	room_ptr.members = append(room_ptr.members,new_player)

	w.WriteHeader(http.StatusOK)
	print_rooms()
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

	print_rooms()
}


// all important debug tools
func print_rooms(){
	m := map[int]interface{}{}

	all_rooms.Range(func(key, value interface{}) bool {
		body := value.(*room)
		room_output := ""

		room_output += "Host: " + body.members[0].name + " || Members:"

		for i := range body.members{
			room_output += " " + body.members[i].name
		}

		// room_output["players"] = body.owner.name + room_mems

		// a, _ := json.Marshal(room_output)

		m[key.(int)] = string(room_output)
		return true
	})

	b, _ := json.MarshalIndent(m, "", " ")
	fmt.Println(string(b))
}
