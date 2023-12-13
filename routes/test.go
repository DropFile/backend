// routes/ping
package routes

import (
	"fmt"
	"strings"

	"cloud/backend/database"

	"github.com/gin-gonic/gin"
)

func HandlePing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World!",
		})
	}
}

func HandleSet(kvStore *database.KVStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data map[string]string
		if err := ctx.BindJSON(&data); err != nil {
			ctx.JSON(400, gin.H{
				"message": "Invalid request body",
			})
			return
		}

		key := data["key"]
		value := data["value"]

		// convert comma separated string to array of string
		values := strings.Split(value, ",")

		if err := kvStore.Set(key, values); err != nil {
			ctx.JSON(500, gin.H{
				"message": fmt.Sprintf("Error setting value %s", err),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"message": "Value set successfully",
		})
	}
}

func HandleGet(kvStore *database.KVStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Query("key")

		result, err := kvStore.Get(key)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": fmt.Sprintf("Error getting value %s", err),
			})
		}

		ctx.JSON(200, gin.H{
			"data":    result,
			"message": "Values get successfully",
		})
	}
}
