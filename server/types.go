package main

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
    id int
    turn int
    hub *Hub
    deck map[card]int
}

type player struct{
    name string
    cards [2]card
    coins int
    client *Client
}
