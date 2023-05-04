package server

import (
	"sync"
	"fmt"
	"encoding/json"
	"strconv"
)

func init_game(lobby *Room){
	// change mode of hub
	draw_cards(lobby)
	lobby.Turn += 1
}

// make proper ids later
var global_id_counter = 0
var id_lock sync.Mutex

func global_id_get() int {
	id_lock.Lock()
	global_id_counter++
	cur_id := global_id_counter
	id_lock.Unlock()

	return cur_id-1
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
