package main

import (
	"encoding/json"
	"fmt"
	//"github.com/bradfitz/slice"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strconv"
)

type config struct {
	Port  int
	Limit int
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
		//fmt.Println("Got message ", msg, " for ", arr)
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
		close(msg.ch)
	}

}

func try_allocate(name string, input []chan message, freeList <-chan int) (int, string) {
	var i, httpStatus int
	var httpMsg string
	select {
	case i = <-freeList:
		output := make(chan bool)
		res := resource{id: i, free: false, owner: name}
		msg := message{data: res, ch: output}
		input[i] <- msg
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

func try_deallocate(id int, input []chan message, freeList <-chan int) (int, string) {
	id--
	var httpStatus int
	var httpMsg string
	output := make(chan bool)
	res := resource{id: id, free: true, owner: ""}
	msg := message{data: res, ch: output}
	input[id] <- msg
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
	fmt.Printf("Loaded config: port %d, limit %d\n", appConfig.Port, appConfig.Limit)

	// Init arr
	arr := aresource{members: []resource{}}
	freeList := make(chan int, appConfig.Limit)
	var input []chan message
	for i := 0; i < appConfig.Limit; i++ {
		r := resource{i + 1, true, ""}
		arr.members = append(arr.members, r)
		ch := make(chan message, 10)
		input = append(input, ch)
		freeList <- i
		go worker(ch, &arr, freeList)
	}

	//slice.Sort(arr.members, arr.Less)
	//fmt.Println("arr: ", arr)

	// Starting gin gonic
	server := gin.Default()

	server.GET("/allocate/:name", func(c *gin.Context) {
		httpStatus, httpMsg := try_allocate(c.Param("name"), input, freeList)
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
			httpStatus, httpMsg = try_deallocate(id, input, freeList)
		}
		c.String(httpStatus, httpMsg)
	})

	server.Run(fmt.Sprintf(":%d", appConfig.Port))
}
