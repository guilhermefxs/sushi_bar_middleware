package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {

	var reply string

	var client, err = rpc.DialHTTP("tcp", "localhost:4040")

	if err != nil {
		log.Fatal("connection error: ", err)
	}

	err = client.Call("API.Restaurant", "", &reply)
	if err != nil {
		return
	}
	fmt.Println("alo")
	fmt.Println(reply)
}
