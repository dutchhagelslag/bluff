package main

import (
	"encoding/json"
    "fmt"
    "time"
	"strconv"
    "math/rand"
)

// assume non malicous players for now
func remove_player(remove_name string, lobby *room){
	// delete room if last player
	if(len(lobby.members) == 1){
		del_room, ok := all_rooms.Load(lobby.id)

		if !ok{
			return
		}

		// for hub, read pump, and write pump
		for i:=0; i<3; i++ {
			del_room.(*room).hub.off <- []byte("off")
		}

		all_rooms.Delete(lobby.id)
		return
	}

	for i := range lobby.members {
		player := lobby.members[i]
		fmt.Println(player.name)


		if(player.name == remove_name){
			// shift everyone
			lobby.members = append(lobby.members[:i], lobby.members[i+1:] ...)
			return
		}
	}
}

// randomly determine cards for players
func draw_cards(lobby *room){
	for i := 0 ; i < len(lobby.members) ; i++ {
		lobby.members[i].cards[0] = get_random_card(lobby.deck)
		lobby.members[i].cards[1] = get_random_card(lobby.deck)
	}
}

// gets card from deck
func get_random_card(deck map[card]int) card{
	var new_card card;
	for {

		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		new_card = card(random.Intn(5) + 1)

		if val, ok := deck[new_card]; (ok && val > 0){
			deck[new_card] -= 1
			return new_card
		}
	}
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

// make a new room
func init_room(host_name string, room_id int) *room{
	members_init := make([]player, 1, 6)

	// init host
	members_init[0] = player{
		name: host_name,
		cards: [2]card{},
		coins: 2,
	}

	new_hub := newHub(room_id)
	go new_hub.run()

	return &room{
		members: members_init,
		turn: 0,
		id: room_id,
		deck: init_deck(),
		hub: new_hub,
	}
}

func init_game(lobby *room){
	// change mode of hub
	draw_cards(lobby)
	lobby.turn += 1
}

// make proper ids later
var global_id_counter = 0

func global_id_get() int {
	global_id_counter++
	return global_id_counter-1
}

func card_to_str(card card) string{
	switch card{
		case Dead:
			return "dead"
		case Duke:
			return "Duke"
		case Captain:
			return "Captain"
		case Ambassador:
			return "Ambassador"
		case Contessa:
			return "Contessa"
		case Assasin:
			return "Assasin"
	}
	return ""
}

// all important debug print
func print_rooms(action string){
	fmt.Println("======================================================================")
	fmt.Println(action)

	m := map[int]interface{}{}

	all_rooms.Range(func(key, value interface{}) bool {
		body := value.(*room)
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
