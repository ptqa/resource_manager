package main

import (
	"encoding/json"
	"fmt"
	//"github.com/bradfitz/slice"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
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

func test() string {
	return "Tests ok\n"
}

func (a *aresource) Less(i, j int) bool {
	if a.members[i].free == true && a.members[j].free == false {
		return true
	}
	return false
}

func (a *aresource) free() (bool, int) {
	free := false
	index := -1
	for i := range a.members {
		if a.members[i].free == true {
			free = true
			index = i
			break
		}
	}
	return free, index
}

func worker(c chan message, arr *aresource) {
	for msg := range c {
		fmt.Println("Got message ", msg, " for ", arr)
		current := arr.members[msg.data.id].free
		if (msg.data.free == false && current == true) || (msg.data.free == true && current == false) {
			arr.members[msg.data.id] = msg.data
			msg.ch <- true
		} else {
			msg.ch <- false
		}
		close(msg.ch)
	}

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
	var input []chan message
	for i := 0; i < appConfig.Limit; i++ {
		r := resource{i + 1, true, "owner1"}
		arr.members = append(arr.members, r)
		ch := make(chan message, 10)
		input = append(input, ch)
		go worker(ch, &arr)
	}

	//slice.Sort(arr.members, arr.Less)
	fmt.Println("arr: ", arr)

	// Starting gin gonic
	server := gin.Default()
	server.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, test())
	})

	server.GET("/allocate/:name", func(c *gin.Context) {
		http_status := http.StatusCreated
		http_msg := "Ops.\n"
		for {
			free, i := arr.free()
			if free == false {
				http_status = http.StatusServiceUnavailable
				http_msg = "Out of resources.\n"
				break
			} else {
				output := make(chan bool)
				res := resource{id: i, free: false, owner: c.Param("name")}
				msg := message{data: res, ch: output}
				input[i] <- msg
				result := <-output
				if result == true {
					http_msg = fmt.Sprintf("r%d\n", i+1)
					break
				}
			}
		}
		c.String(http_status, http_msg)
	})

	server.Run(fmt.Sprintf(":%d", appConfig.Port))
}
