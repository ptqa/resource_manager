package resource_manager

import (
	//"log"
	//"github.com/bradfitz/slice"
	"fmt"
	"github.com/dgryski/go-jump"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Resource struct {
	id    int
	free  bool
	owner string
}

type Message struct {
	data Resource
	ch   chan bool
}

type Resources struct {
	members  []Resource
	input    []chan Message
	freeList chan int
}

/*
func (a *aresource) Less(i, j int) bool {
	if a.members[i].free == true && a.members[j].free == false {
		return true
	}
	return false
}
*/

// Simple and fast hashring
func choose_worker(i int, n int) int {
	i64 := uint64(i)
	place := jump.Hash(i64, n)
	return int(place)
}

func (a *Resources) worker(c <-chan Message) {
	for msg := range c {
		if len(a.members) < msg.data.id || msg.data.id < 0 {
			msg.ch <- false
		} else {
			current := a.members[msg.data.id].free
			if msg.data.free != current {
				a.members[msg.data.id] = msg.data
				if msg.data.free == true {
					a.freeList <- msg.data.id
				}
				msg.ch <- true
			} else {
				msg.ch <- false
			}
		}
		close(msg.ch)
	}
}

func (a *Resources) try_allocate(name string, workers int) (int, string) {
	var i, httpStatus int
	var httpMsg string
	select {
	case i = <-a.freeList:
		output := make(chan bool)
		res := Resource{id: i, free: false, owner: name}
		msg := Message{data: res, ch: output}
		place := choose_worker(i, workers)
		a.input[place] <- msg
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

func (a *Resources) try_deallocate(id int, workers int) (int, string) {
	id--
	var httpStatus int
	var httpMsg string
	output := make(chan bool)
	res := Resource{id: id, free: true, owner: ""}
	msg := Message{data: res, ch: output}
	place := choose_worker(id, workers)
	a.input[place] <- msg
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

func (a *Resources) list() string {
	return fmt.Sprintf("%s", a.members)
}

func main() {

	// Parsing config
	var appConfig Config
	err := appConfig.Load("config.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded config: port %d, limit %d, workers %d\n", appConfig.Port, appConfig.Limit, appConfig.Workers)

	// Init arr
	arr := Resources{}

	arr.freeList = make(chan int, appConfig.Limit)

	for i := 0; i < appConfig.Limit; i++ {
		r := Resource{i + 1, true, ""}
		arr.freeList <- i
		arr.members = append(arr.members, r)
	}

	for i := 0; i < appConfig.Workers; i++ {
		ch := make(chan Message, 10)
		arr.input = append(arr.input, ch)
		go arr.worker(ch)
	}

	// Starting gin gonic
	server := gin.Default()

	server.GET("/allocate/:name", func(c *gin.Context) {
		httpStatus, httpMsg := arr.try_allocate(c.Param("name"), appConfig.Workers)
		c.String(httpStatus, httpMsg)
	})

	server.GET("/deallocate/r:id", func(c *gin.Context) {
		var httpStatus int
		var httpMsg string
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil && id <= len(arr.input) {
			httpMsg = "Not allocated\n"
			httpStatus = 404
		} else {
			httpStatus, httpMsg = arr.try_deallocate(id, appConfig.Workers)
		}
		c.String(httpStatus, httpMsg)
	})

	server.GET("/list", func(c *gin.Context) {
		httpStatus := 200
		httpMsg := arr.list()
		c.String(httpStatus, httpMsg)
	})

	server.Run(fmt.Sprintf(":%d", appConfig.Port))
}
