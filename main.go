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
	state bool
	owner string
}

type aresource struct {
	members []resource
	free    int
}

func test() string {
	return "Tests ok\n"
}

func (a *aresource) Less(i, j int) bool {
	if a.members[i].state == true && a.members[j].state == false {
		return true
	}
	return false
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
	arr := aresource{members: []resource{}, free: 0}
	for i := 0; i < appConfig.Limit; i++ {
		r := resource{i + 1, false, "owner1"}
		arr.members = append(arr.members, r)
		arr.free++
	}

	//slice.Sort(arr.members, arr.Less)
	fmt.Println("arr: ", arr)

	// Starting gin gonic
	server := gin.Default()
	server.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, test())
	})

	server.GET("/allocate/:name", func(c *gin.Context) {
		http_status := http.StatusOK
		http_msg := "OK\n"
		if arr.free <= 0 {
			http_status = http.StatusServiceUnavailable
			http_msg = "Out of resources.\n"
		}
		c.String(http_status, http_msg)
	})

	/*
		server.GET("/data/", func(c *gin.Context) {
			data := <-input
			c.String(http.StatusOK,
			"Id: %d, Name: %s, Data: %s \n",
			data.id,
			data.name,
			data.data)
		})

		server.PUT("/data/:id/:name/:data", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.String(http.StatusInternalServerError, "Can't convert id to integer\n")
			}
			go func() {
				input <- Data{id: id, name: c.Param("name"), data: c.Param("data")}
			}()
			c.String(http.StatusOK, "OK!")
		})
	*/
	server.Run(fmt.Sprintf(":%d", appConfig.Port))
}
