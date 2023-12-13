package main

import (
	// "net/http"
	"cloud/backend/database"
	"cloud/backend/routes"
	"github.com/gin-gonic/gin"
	"log"
)

const dbPath string = "../../database"

func main() {

	// database connection
	kvStore, err := database.NewKVStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer kvStore.Close()

	router := gin.Default()

	router.GET("/ping", routes.HandlePing())
	router.POST("/db/test", routes.HandleSet(kvStore))
	router.GET("/db/test", routes.HandleGet(kvStore))
	router.POST("/file/upload", routes.HandleUpload(kvStore))

	router.Run()
}
