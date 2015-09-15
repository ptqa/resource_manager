package main

import (
	//"log"
	//"github.com/bradfitz/slice"
	"fmt"
	"github.com/dgryski/go-jump"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Message struct {
	data Resource
	ch   chan bool
}

// Simple and fast hashring
func choose_worker(i int, n int) int {
	i64 := uint64(i)
	place := jump.Hash(i64, n)
	return int(place)
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
