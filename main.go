package main

import (
	// "net/http"
	"cloud/backend/database"
	"cloud/backend/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	dbPath := "../../database"

	// database connection
	kvStore, err := database.NewKVStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer kvStore.Close()

	router := gin.Default()

	router.GET("/ping", routes.HandlePing())
	router.POST("/test", routes.HandleSet(kvStore))
	router.GET("/test", routes.HandleGet(kvStore))

	router.Run()
}
