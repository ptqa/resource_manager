package main

import (
	//"log"
	//"github.com/bradfitz/slice"
	"fmt"
	//"github.com/davecheney/profile"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Message struct {
	data Resource
	ch   chan bool
}

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()

	// Parsing config
	var appConfig Config
	err := appConfig.Load("config.json")
	if err != nil {
		panic(err)
	}

	// Init arr
	var arr Resources
	arr.Init(appConfig)

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

	server.GET("/list/:name", func(c *gin.Context) {
		httpStatus := 200
		httpMsg := arr.search(c.Param("name"))
		c.String(httpStatus, httpMsg)
	})

	server.Run(fmt.Sprintf(":%d", appConfig.Port))
}
