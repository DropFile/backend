package main

import (
	// "net/http"
	"github.com/gin-gonic/gin"
	"cloud/backend/routes"
)

func main() {
	router := gin.Default()
	router.GET("/ping", routes.Ping)

	router.Run()
}