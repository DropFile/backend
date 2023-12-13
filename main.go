package main

import (
	// "net/http"
	"cloud/backend/database"
	"cloud/backend/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	router.StaticFS("/file/fetch", http.Dir(fmt.Sprintf("%s/", routes.Filestorage)))

	router.GET("/ping", routes.HandlePing())
	router.POST("/db/test", routes.HandleSet(kvStore))
	router.GET("/db/test", routes.HandleGet(kvStore))
	router.POST("/file", routes.HandleUpload(kvStore))
	router.GET("/file", routes.HandleFileMetadata(kvStore))

	router.Run()
}
