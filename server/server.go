package main

import (
    "strconv"
    "fmt"
    // "os"
    "net/http"
    "log"
    // "time"
    "regexp"
    "errors"
	"sync"
	"bluff/server/spinner"
	"encoding/json"
    "github.com/julienschmidt/httprouter"
)

// assume non malicous players for now
type id    string
type card int8

const(
    Duke card = iota
    Assasin
    Contessa
    Captain
    Ambassador
)

type player struct{
    name string
    cards [2]card
    coins int
}

type game_state struct{
    turn int
    deck map[card]int 
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

type room struct{
    owner player
    members []player
    cur_game_state game_state
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

func init_room(host_name string) *room{
	host := player{
		name: host_name,
		cards: [2]card{},
		coins: 2,
	}
	
	start_game_state := game_state{
		turn: 0,
		deck: init_deck(),
	}

	return &room{
		owner: host,
		members: init_members(host),
		cur_game_state: start_game_state,
	}
}

func init_members(host player) []player{
	mems := make([]player, 1, 6)
	mems[0] = host
	return mems
}

// func (p *Page) save() error{
//     filename := p.Title + ".txt"
//     return os.WriteFile(filename, p.Body, 0600)
// }

// func loadPage(title string) (*Page, error){
//     filename := title + ".txt"
//     body, err := os.ReadFile(filename)
//     if err != nil {
//         return nil, err
//     }

//     return &Page{Title: title, Body: body}, nil
// }

// Tax :Duke

// Bluff block foreign aid :Duke
// Bluff block assasin :Contessa
// Bluff block stealing :Captain :Ambassador

// return all card info as json
func card_info_handler(w http.ResponseWriter, r *http.Request){
	info := map[string][2]string{
		"Duke": [2]string{"Take 3 coins from the Treasury",
			              "Block Foreign Aid"},

		"Contessa": [2]string{"",
			              "Block Assasination"},


		// "Assasin":"Pay 3 coins to the Treasury and launch an assassination against another player. If successful that player immediately loses an influence. (Can be blocked by the Contessa)",

		// "Contessa":"",

		// "Captain":"Take 2 coins from another player. If they only have one coin, take only one. (Can be blocked by the Ambassador or the Captain)",

		// "Ambassador":"Exchange cards with the Court. First take 2 random cards from the Court deck. Choose which, if any, to exchange with your face-down cards. Then return two cards to the Court deck.",


	}

	json_info, _ := json.Marshal(info)
	w.Write(json_info)
}

// func get_card_actions(card card) (){
//     switch card{
//         case Duke:
//             return []

//         case Assasin:
//             return []

//         case Contessa:
//             return []

//         case Captain:
//             return []

//         case Ambassador:
//             return []
//     }
// }



// server to dos:
// implement tls
// set addr
//

// session manager

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


// all important debug tools
func print_rooms(){
	m := map[int]interface{}{}

	all_rooms.Range(func(key, value interface{}) bool {
		body := value.(*room)
		room_output := ""

		room_output += "Host: " + body.owner.name + " || Members:" 

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

// func print_hand(){

// }

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
	id, _ := strconv.Atoi(ps.ByName("room_id"))
	room_ptr, ok := all_rooms.Load(id)

	// room doesn't exist
	if(!ok){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ok = insert_player(room_ptr.(*room), ps.ByName("player_name"))	

	// room is full
	if(!ok){
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func insert_player(room_ptr *room, player_name string) bool{
	// limit to capacity
	if(len(room_ptr.members) == cap(room_ptr.members)){
		return false
	}

	new_player := player{
		name: player_name,
		cards: [2]card{},
		coins: 2,
	}

	room_ptr.members = append(room_ptr.members,new_player)
	return true
}
	
func start_game(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	
}
