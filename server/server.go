package main

import (
	"sync"
	"fmt"
	"encoding/json"
	"strconv"
)

func init_game(lobby *Room){
	// change mode of hub
	draw_cards(lobby)
	lobby.turn += 1
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

		room_output += "Host: " + body.members[0].name + " || Members:"
		for i := range body.members{
			room_output += " " + body.members[i].name
		}
		for i := range body.members{
			hands += " " + body.members[i].name + ": "
			hands += card_to_str(body.members[i].cards[0]) + " and "
			hands += card_to_str(body.members[i].cards[1]) + "||"
		}

		// room_output["players"] = body.owner.name + room_mems

		// a, _ := json.Marshal(room_output)

		deck, _ := json.MarshalIndent(body.deck, "", " ")
		fmt.Println(string(deck))

		m[key.(int)] = string(room_output + hands)
		fmt.Println("LobbyLen: " + strconv.Itoa(len(body.members)))
		fmt.Println("LobbyCap: " + strconv.Itoa(cap(body.members)))

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
