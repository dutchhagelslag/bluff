package main

import (
    "fmt"
    // "os"
    "net/http"
    "log"
    // "time"
    "regexp"
    "errors"
	// "sync"
	"bluff/server/spinner"
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
    players []player
    deck []card // 3 of each cards - cards in players hands (ignore till ambassador impl)
}

type lobby struct{
    owner player
    members []player
    size int
    cur_game_state game_state
}

const max_rooms = 5
var all_lobbies [max_rooms]lobby

// add error handling

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

// return all card info as json
func card_info_handler(w http.ResponseWriter, r *http.Request){

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



func main(){
    health_check := make(chan string)

    go spinner.Spin(health_check)

    router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)

    log.Fatal(http.ListenAndServe(":8080", router))
}


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

