package server

import(
	"fmt"
	// "strconv"
	// "bytes"
	"errors"
	"sync"
	"encoding/json"
)

type Room struct{
   Members []*Player
   State string 	// lobby, ingame, get-action, get-bluff-call
   Turn int
   Id int                    `json:"id"`

   Deck map[Card]int         `json:"-"`
   Join chan *Player         `json:"-"`

   Receive chan []byte       `json:"-"`
   Members_mut sync.Mutex    `json:"-"`
}


func init_room(room_id int) *Room{
	new_room := Room{
		Members: make([]*Player, 0, 6),
		Turn: 0,
		Id: room_id,
		Deck: init_deck(),
		Join:  make(chan *Player),
		Receive: make(chan []byte),
		State: "lobby",
		Members_mut: sync.Mutex{},
	}

	go new_room.run()

	all_rooms.Store(room_id, &new_room)

	return &new_room
}

// return false if already exists or capacity met
func (room *Room) add(player *Player) error{
	if len(room.Members) == cap(room.Members){
		return errors.New("Error: room is full")
	}

	if room.Turn != 0{
		return errors.New("Error: Game has already started")
	}

	for i := range(room.Members){
		if room.Members[i].Name == player.Name {
			return errors.New("Error: Player Name taken")
		}
	}

	room.Members_mut.Lock()
	room.Members = append(room.Members,player)
	room.Members_mut.Unlock()

	player.Room = room
	return nil
}

// assume non malicous players for now
func (room *Room) remove(remove_Name string){

	for i := range room.Members {
		player := room.Members[i]

		if(player.Name == remove_Name){
			room.Members_mut.Lock()

			player.Conn.Close()

			// DisConnect and delete Room if last player
			if(len(room.Members) == 1){
				all_rooms.Delete(room.Id)
				return
			}


			// shift everyone
			room.Members = append(room.Members[:i], room.Members[i+1:] ...)

			room.Members_mut.Unlock()
			return
		}
	}
}

func (room *Room) get_player(Name string) *Player{
	for _, member := range room.Members {
		if member.Name == Name{
			return member
		}
	}
	return nil
}


// runs lobby then game logic
func (room *Room) run() {
	for {
		room.run_lobby()

		if winner := room.run_game(); winner != nil{
			room.announce_state()
		}else{
			room.announce_state()
		}
	}
}

// lobby stage ========================

func (room *Room) run_lobby() {
	for {
		if(room.is_room_rdy()){
			return
		}

		room.ready_up(<-room.Receive)
		room.announce_state()
	}
}

func (room *Room) is_room_rdy() bool{
	for _, member := range room.Members{
		if !member.Rdy{
			return false
		}
	}

	// I wont let you play alone
	if len(room.Members) == 1 {
		return false
	}

	return true
}

func (Room *Room) ready_up(player_Name []byte){
	Name := string(player_Name)

	if Name[0] == '!'{
		player := Room.get_player(Name[1:])
		if player == nil{
			return
		}

		player.Rdy = false
	}else{
		player := Room.get_player(Name)
		if player == nil{
			return
		}

		player.Rdy = true
	}
}

// ===================================




// game stage ========================

func (Room *Room) announce_state(){
	Room_json, err := json.Marshal(Room)

	if err != nil{
		fmt.Println("Invalid game state")
	}

	for _, member := range Room.Members {
		member.Send<-Room_json
	}
}

func (room *Room) run_game() *Player{

	for {
		num_alive := room.players_left()

		// when the winner dc before annoucement loop -> lobby
		if(num_alive == 0){
			return nil
		}

		if(num_alive == 1){
			for _, member := range(room.Members){
				if member.is_alive(){
					return member
				}
			}
		}

		cur_player := room.whose_turn()

		if cur_player == nil {
			return nil
		}

		// announce whose turn it is
		room.announce_state()

		// ask action player to choose an action
		// player_action := room.ask("get-action", cur_player)

		// room.announce("") // action of cur_player

		// room.announce("Challenge:" + cur_player.Name)

		// challenger := room.countdown_ask(20, nil)

		// if challenger != nil{
		// 	room.challenge(cur_player, challenger)
		// }else{
		// 	room.process_action(cur_player, player_action)
		// }

		room.Turn++
	}

}

func (room *Room) ask(action string, player *Player) Action{
	ask_json, _ := json.Marshal(
		struct{
			State string
		} {
			action,
		},
	)
	player.Send<-ask_json

	// resp := <- player.Send



	return Action(3)
}


func (room *Room) players_left() int{
	alive := 0
	for _, member := range(room.Members){
		if member.is_alive(){
			alive++
		}
	}
	return alive
}

func (room *Room) whose_turn() *Player{
	for{
		if room.players_left() == 0 {
			return nil
		}

		index := room.Turn % len(room.Members)
		room.Turn++

		if(room.Members[index].is_alive()){
			return room.Members[index]
		}
	}
}
