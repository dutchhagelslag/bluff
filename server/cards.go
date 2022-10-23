package main

import (
	// "encoding/json"
    // "fmt"
    "time"
	// "strconv"
    "math/rand"
)

type Id    string
type Card int8

const(
    Dead Card = iota
    Duke
    Assasin
    Contessa
    Captain
    Ambassador
)

// randomly determine cards for players
func draw_cards(room *Room){
	for i := 0 ; i < len(room.members) ; i++ {
		room.members[i].cards[0] = get_random_card(room.deck)
		room.members[i].cards[1] = get_random_card(room.deck)
	}
}

// gets card from deck
func get_random_card(deck map[Card]int) Card{
	var new_card Card;
	for {

		seed := rand.NewSource(time.Now().UnixNano())
		random := rand.New(seed)
		new_card = Card(random.Intn(5) + 1)

		if val, ok := deck[new_card]; (ok && val > 0){
			deck[new_card] -= 1
			return new_card
		}
	}
}

// 3 of each cards - cards in players hands (ignore till ambassador impl)
func init_deck() map[Card]int{
	return map[Card]int{
		Duke:3,
		Assasin:3,
		Contessa:3,
		Captain:3,
		Ambassador:3,
	}
}

func card_to_str(card Card) string{
	switch card{
		case Dead:
			return "Dead"
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
