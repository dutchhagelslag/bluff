package main

import(
	// "fmt"
	// "strconv"
	// "bytes"
	"errors"
	"sync"
	"encoding/json"
)

type Room struct{
    members []*Player
    id int
    turn int
    deck map[Card]int

	join chan *Player

	// recieve chan PlayerResp
	receive chan []byte

	members_mut sync.Mutex

	// lobby, ingame, get-action, get-bluff-call
	state string
}

// make a new room
func init_room(room_id int) *Room{
	new_room := Room{
		members: make([]*Player, 0, 6),
		turn: 0,
		id: room_id,
		deck: init_deck(),
		join:  make(chan *Player),
		receive: make(chan []byte),
		state: "lobby",
		members_mut: sync.Mutex{},
	}

	go new_room.run()

	all_rooms.Store(room_id, &new_room)

	return &new_room
}

// return false if already exists or capacity met
func (room *Room) add(player *Player) error{
	if len(room.members) == cap(room.members){
		return errors.New("Error: Room is full")
	}

	if room.turn != 0{
		return errors.New("Error: Game has already started")
	}

	for i := range(room.members){
		if room.members[i].name == player.name {
			return errors.New("Error: Player name taken")
		}
	}

	room.members_mut.Lock()
	room.members = append(room.members,player)
	room.members_mut.Unlock()

	player.room = room
	return nil
}

// assume non malicous players for now
func (room *Room) remove(remove_name string){

	for i := range room.members {
		player := room.members[i]

		if(player.name == remove_name){
			room.members_mut.Lock()

			player.conn.Close()

			// Disconnect and delete room if last player
			if(len(room.members) == 1){
				all_rooms.Delete(room.id)
				return
			}


			// shift everyone
			room.members = append(room.members[:i], room.members[i+1:] ...)

			room.members_mut.Unlock()
			return
		}
	}
}

func (room *Room) get_player(name string) *Player{
	for _, member := range room.members {
		if member.name == name{
			return member
		}
	}
	return nil
}


// func (room *Room) broadcast_state(){
// 	for player := range room.members{
// 		player.send <- game_state
// 	}
// }

// runs lobby then game logic
func (room *Room) run() {
	for {
		room.run_lobby()

		if winner := room.run_game(); winner != nil{
			room.announce("winner: " + winner.name)
		}else{
			room.announce("winner: unknown")
		}
	}
}

// lobby stage ========================

func (room *Room) run_lobby() {
	for {
		if(room.is_room_rdy()){
			room.announce("Game Start")
			return
		}

		room.ready_up(<-room.receive)
		room.announce()
	}
}

func (room *Room) is_room_rdy() bool{
	for _, member := range room.members{
		if !member.rdy{
			return false
		}
	}

	// I wont let you play alone
	if len(room.members) == 1 {
		return false
	}

	return true
}

func (room *Room) ready_up(player_name []byte){
	name := string(player_name)

	if name[0] == '!'{
		room.get_player(name[1:]).rdy = false
	}else{
		room.get_player(name).rdy = true
	}
}

// ===================================




// game stage ========================

func (room *Room) announce(msg []byte){
	for i := range room.members {
		room.members[i].send <- msg
	}
}

func (room *Room) run_game() *Player{

// 	for {
// 		num_alive := room.run_alive()

// 		// when the winner dc before annoucement loop -> lobby
// 		if(num_alive == 0){
// 			return nil
// 		}

// 		if(num_alive == 1){
// 			for i := range(room.members){
// 				if room.members[i].is_alive(){
// 					return room.members[i].name
// 				}
// 			}
// 		}

// 		cur_player := room.whose_turn()

// 		if cur_player == nil {
// 			return nil
// 		}

// 		room.announce("Turn:" + cur_player.name)

// 		player_action := room.countdown_ask(20, cur_player)

// 		room.announce("") // action of cur_player

// 		room.announce("Challenge:" + cur_player.name)

// 		challenger := room.countdown_ask(20, nil)

// 		if challenger != nil{
// 			room.challenge(cur_player, challenger)
// 		}else{
// 			room.process_action(cur_player, player_action)
// 		}

// 		room.turn++
// 	}

// }

// func (room *Room) countdown(secs int){
// 	// timer and ping everyone for each second
// }

// func (room *Room) find_challenger(){
// 	start_index := room.turn

// 	for i := range(room.members){
// 		cur_index = (start_index + i) % len(room.members)
// 	}
	return nil
}


// func (room *Room) process_action(){

// }

// func (room *Room) challenge(){

// }

// func (room *Room) whose_turn() *Player{
// 	for{
// 		if room.num_alive == 0 {
// 			return nil
// 		}

// 		index = room.turn % len(room.members)

// 		if(room.member[index].is_alive()){
// 			return room.member[index]
// 		}

// 		room.turn++
// 	}
// }


// func (room *Room) num_alive() int{
// 	num_alive := 0

// 	for i := range(room.members){
// 		if room.members[i].is_alive(){
// 			num_alive++
// 		}
// 	}

// 	return num_alive
// }

// ===================================
