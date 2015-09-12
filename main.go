package main

import (
	"encoding/json"
	"fmt"
	//"github.com/bradfitz/slice"
	"github.com/dgryski/go-jump"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strconv"
)

type config struct {
	Port    int
	Limit   int
	Workers int
}

type resource struct {
	id    int
	free  bool
	owner string
}

type message struct {
	data resource
	ch   chan bool
}

type aresource struct {
	members []resource
}

/*
func (a *aresource) Less(i, j int) bool {
	if a.members[i].free == true && a.members[j].free == false {
		return true
	}
	return false
}
*/

func worker(c <-chan message, arr *aresource, freeList chan<- int) {
	for msg := range c {
		if len(arr.members) < msg.data.id || msg.data.id < 0 {
			msg.ch <- false
		} else {
			current := arr.members[msg.data.id].free
			if msg.data.free != current {
				arr.members[msg.data.id] = msg.data
				if msg.data.free == true {
					freeList <- msg.data.id
				}
				msg.ch <- true
			} else {
				msg.ch <- false
			}
		}
		close(msg.ch)
	}
}

func choose_worker(i int, n int) int {
	i64 := uint64(i)
	place := jump.Hash(i64, n)
	return int(place)
}
func try_allocate(name string, input []chan message, freeList <-chan int, workers int) (int, string) {
	var i, httpStatus int
	var httpMsg string
	select {
	case i = <-freeList:
		output := make(chan bool)
		res := resource{id: i, free: false, owner: name}
		msg := message{data: res, ch: output}
		place := choose_worker(i, workers)
		input[place] <- msg
		result := <-output
		if result == true {
			httpStatus = 200
			httpMsg = fmt.Sprintf("r%d\n", i+1)
		}
	default:
		httpStatus = 503
		httpMsg = "Out of resources.\n"
	}
	return httpStatus, httpMsg
}

func try_deallocate(id int, input []chan message, freeList <-chan int, workers int) (int, string) {
	id--
	var httpStatus int
	var httpMsg string
	output := make(chan bool)
	res := resource{id: id, free: true, owner: ""}
	msg := message{data: res, ch: output}
	place := choose_worker(id, workers)
	input[place] <- msg
	result := <-output
	if result == true {
		httpStatus = 204
		httpMsg = ""
	} else {
		httpMsg = "Not allocated\n"
		httpStatus = 404
	}
	return httpStatus, httpMsg
}

func main() {

	// Parsing config
	var appConfig config
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(configFile, &appConfig); err != nil {
		log.Fatal(err)
	}
	if appConfig.Port < 1 || appConfig.Limit < 1 || appConfig.Workers < 2 {
		panic(fmt.Sprintf("Invalid config\n"))
	}
	fmt.Printf("Loaded config: port %d, limit %d, workers %d\n", appConfig.Port, appConfig.Limit, appConfig.Workers)

	// Init arr
	arr := aresource{members: []resource{}}

	freeList := make(chan int, appConfig.Limit)
	var input []chan message

	for i := 0; i < appConfig.Limit; i++ {
		r := resource{i + 1, true, ""}
		freeList <- i
		arr.members = append(arr.members, r)
	}

	for i := 0; i < appConfig.Workers; i++ {
		ch := make(chan message, 10)
		input = append(input, ch)
		go worker(ch, &arr, freeList)
	}

	// Starting gin gonic
	server := gin.Default()

	server.GET("/allocate/:name", func(c *gin.Context) {
		httpStatus, httpMsg := try_allocate(c.Param("name"), input, freeList, appConfig.Workers)
		c.String(httpStatus, httpMsg)
	})

	server.GET("/deallocate/r:id", func(c *gin.Context) {
		var httpStatus int
		var httpMsg string
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil && id <= len(input) {
			httpMsg = "Not allocated\n"
			httpStatus = 404
		} else {
			httpStatus, httpMsg = try_deallocate(id, input, freeList, appConfig.Workers)
		}
		c.String(httpStatus, httpMsg)
	})

	server.Run(fmt.Sprintf(":%d", appConfig.Port))
}
