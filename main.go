// This file contains main function
// and runs http server with gin
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()

	// Parsing config
	var appConfig Config
	err := appConfig.Load("config.json")
	if err != nil {
		panic(err)
	}

	// Init arr
	var resources Resources
	resources.Init(appConfig)

	// Starting gin gonic
	server := gin.Default()

	server.GET("/allocate/:name", func(c *gin.Context) {
		httpStatus := 200
		var httpMsg string
		id, err := resources.tryAllocate(c.Param("name"), appConfig.Workers)
		if err != nil {
			httpMsg = "Out of resources.\n"
			httpStatus = 503
		}
		httpMsg = "r" + strconv.Itoa(id) + "\n"
		c.String(httpStatus, httpMsg)
	})

	server.GET("/deallocate/r:id", func(c *gin.Context) {
		var httpStatus int
		var httpMsg string
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil && id <= len(resources.input) {
			httpMsg = "Not allocated\n"
			httpStatus = 404
		} else {
			err := resources.tryDeallocate(id, appConfig.Workers)
			if err == nil {
				httpMsg = ""
				httpStatus = 204
			}
		}
		c.String(httpStatus, httpMsg)
	})

	server.GET("/list", func(c *gin.Context) {
		httpStatus := 200
		httpMsg := resources.List()
		c.String(httpStatus, httpMsg)
	})

	server.GET("/list/:name", func(c *gin.Context) {
		httpStatus := 200
		httpMsg := resources.Search(c.Param("name"))
		c.String(httpStatus, httpMsg)
	})

	server.GET("/reset", func(c *gin.Context) {
		httpStatus := 204
		httpMsg := ""
		resources.Reset(appConfig.Workers)
		c.String(httpStatus, httpMsg)
	})

	server.NoRoute(func(c *gin.Context) {
		httpStatus := 400
		httpMsg := "Bad request.\n"
		c.String(httpStatus, httpMsg)
	})

	server.Run(fmt.Sprintf(":%d", appConfig.Port))
}
