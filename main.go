package main

import (
	"encoding/json"
	"fmt"
	"github.com/bradfitz/slice"
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
	arr := aresource{members: []resource{}}
	for i := 0; i < appConfig.Limit; i++ {
		r := resource{i, false, "owner1"}
		arr.members = append(arr.members, r)
	}

	//fmt.Println("Allocated arr size:", len(arr))
	arr.members[2] = resource{99, true, "wut"}
	arr.members[1] = resource{0, true, "wut2"}
	slice.Sort(arr.members, arr.Less)
	fmt.Println("arr: ", arr)

	server := gin.Default()
	server.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, test())
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
