package main

import (
	// "net/http"
	"cloud/backend/routes"
	"fmt"
	"net/http"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// router.Use(addHeaders)

	// config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	// router.Use(cors.New(config))
	router.StaticFS("/file/fetch", http.Dir(fmt.Sprintf("%s/", routes.Filestorage)))

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	router.POST("/file", routes.HandleUpload())
	router.GET("/file", routes.HandleFileMetadata())

	router.Run()
}

