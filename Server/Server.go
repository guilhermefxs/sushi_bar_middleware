package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

type API int

type RestaurantQueue struct {
	Queue           sync.Mutex
	eating, waiting int
	must_wait       bool
}

var r = RestaurantQueue{
	Queue:     sync.Mutex{},
	eating:    0,
	waiting:   0,
	must_wait: false,
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Eating() {
	time.Sleep(5 * time.Second)
	GoOutTheRestaurant()
}

func GoOutTheRestaurant() {
	r.Queue.Lock()
	defer r.Queue.Unlock()
	r.eating -= 1
	fmt.Println("dando o fora")
	fmt.Printf("agora tem %d comendo \n", r.eating)
	if r.eating == 0 {
		fmt.Println("todos sairam do restaurante")
		n := min(r.waiting, 5)
		r.waiting -= n
		r.eating += n
		fmt.Printf("%d entratam no restaurante\n", r.eating)
		for i := 0; i < r.eating; i++ {
			go Eating()
		}
		r.must_wait = r.eating == 5
		if r.must_wait {
			fmt.Println("restaurante cheio")
		}
	}
}

func EnterRestaurant() bool {
	r.Queue.Lock()
	defer r.Queue.Unlock()
	if r.must_wait {
		r.waiting += 1
		fmt.Println("entrando na fila do restaurante")
		return false
	} else {
		r.eating += 1
		r.must_wait = r.eating == 5
		fmt.Println("entrando na mesa do restaurante")
		if r.must_wait {
			fmt.Println("restaurante cheio")
		}
		return true
	}
}

func (a *API) Restaurant(nothing string, reply *string) error {
	isInsideRestaurant := EnterRestaurant()
	if isInsideRestaurant {
		*reply = "Client entered the restaurant"
		go Eating()
	} else {
		*reply = "Client is waiting to get in the restaurant"
	}
	return nil
}

func main() {
	var api = new(API)
	err := rpc.Register(api)

	if err != nil {
		log.Fatal("error registering API", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4040")

	if err != nil {
		log.Fatal("Listener error", err)
	}

	log.Printf("serving rpc on port %d", 4040)
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("error serving: ", err)
	}

}
