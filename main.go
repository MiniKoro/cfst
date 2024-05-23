package main

import (
	"cfst/core"
	"cfst/scheduling"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	core.LoadYmlConfig()
	scheduling.CronInit()
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/", func(context *gin.Context) {
		json := scheduling.ReadResult()
		context.String(http.StatusOK, json)
	})
	err := router.Run(":9846")
	if err != nil {
		return
	}
}
