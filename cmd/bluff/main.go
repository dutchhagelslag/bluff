package main

import bluff "github.com/dutchhagelslag/bluff/internal/server"

func main(){
	// health_check := make(chan string)
	// go spin(health_check)
	bluff.RunServer()
}
