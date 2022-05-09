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
	QueueMutex      sync.Mutex
	Queue           []string
	eating, waiting int
	must_wait       bool
}

var r = RestaurantQueue{
	QueueMutex: sync.Mutex{},
	Queue:      make([]string, 0),
	eating:     0,
	waiting:    0,
	must_wait:  false,
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Eating(name string) {
	time.Sleep(5 * time.Second)
	GoOutTheRestaurant(name)
}

func GoOutTheRestaurant(name string) {
	r.QueueMutex.Lock()
	defer r.QueueMutex.Unlock()
	r.eating -= 1
	fmt.Println(name + " saindo do restaurante")
	fmt.Printf("agora tem %d comendo \n", r.eating)
	if r.eating == 0 {
		fmt.Println("todos sairam do restaurante")
		n := min(r.waiting, 5)
		r.waiting -= n
		r.eating += n
		for i := 0; i < r.eating; i++ {
			fmt.Printf("%s entrou no restaurante\n", r.Queue[0])
			go Eating(r.Queue[0])
			r.Queue = r.Queue[1:]
		}
		fmt.Printf("%d entratam no restaurante\n", r.eating)
		r.must_wait = r.eating == 5
		if r.must_wait {
			fmt.Println("restaurante cheio")
		}
	}
}

func TryToEnterRestaurant(name string) bool {
	r.QueueMutex.Lock()
	defer r.QueueMutex.Unlock()
	if r.must_wait {
		r.waiting += 1
		r.Queue = append(r.Queue, name)
		fmt.Println(name + " entrando na fila do restaurante")
		return false
	} else {
		r.eating += 1
		r.must_wait = r.eating == 5
		fmt.Println(name + " entrando na mesa do restaurante")
		if r.must_wait {
			fmt.Println("restaurante cheio")
		}
		return true
	}
}

func (a *API) Restaurant(name string, reply *string) error {
	isInsideRestaurant := TryToEnterRestaurant(name)
	if isInsideRestaurant {
		*reply = name + " entered the restaurant"
		go Eating(name)
	} else {
		*reply = name + " is waiting to get in the restaurant"
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
