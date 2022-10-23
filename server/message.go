package main

import (
	"fmt"
	"encoding/json"
)

type ReqAction int8
const(
	PlayerTurn ReqAction = iota
	ChallengerTurn
	Announce
	DisplayWinner
)

type PlayerResp int8
const(
	Pass PlayerResp = iota   // +1 coin
	Buyout                   // -1 card for target

	CashOut                  // +2 coin & begin challenge
	Steal                    // +2 coin, -2 coins of target & challenge
	Assasinate  			 // -1 card for target & challenge

	Challenge
)


type BluffMessage struct{
	action string
	payload any
}

type RoomReq struct{
	members []*Player
	action ReqAction
	msg string
}

func (room *Room) getRoomState() RoomState{

}


type LobbyState struct{
	room_id int
	members []string
}

func (room *Room) getLobbyState() LobbyState{
	members_list := make([]string)

	for i := range(room.members){
		members_list.append(room.members[i].name)
	}

	return LobbyState{
		room_id: room.id,
		members: members_list,
	}
}

/*

room -> player
Turn:
Challenge:


*/
