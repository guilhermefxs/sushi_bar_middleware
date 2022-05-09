package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

func main() {
	//o nome de quem vai entrar no restaurante eh dada como argumento an cli
	argsWithoutProg := os.Args[1:]

	var reply string

	var client, err = rpc.DialHTTP("tcp", "localhost:4040")

	if err != nil {
		log.Fatal("connection error: ", err)
	}

	err = client.Call("API.Restaurant", argsWithoutProg[0], &reply)
	if err != nil {
		return
	}
	fmt.Println(reply)
}
