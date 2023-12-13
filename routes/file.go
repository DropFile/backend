package routes

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"cloud/backend/database"
	"cloud/backend/utils"

	"github.com/gin-gonic/gin"
)

const filestorage string = "../../filestorage"

func HandleUpload(kvStore *database.KVStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, err := utils.GenerateRandomString(20)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error generating Random String",
			})
			return
		}

		fileStoragePath := fmt.Sprintf("%s/%s", filestorage, key)
		err = utils.CreateFolder(fileStoragePath)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error creating folder",
			})
			return
		}

		var fileNames []string
		numberOfFiles, err := strconv.Atoi(ctx.Request.FormValue("numberOfFiles"))
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "Error converting to Integer",
			})
			return
		}

		var wg sync.WaitGroup

		for num := 0; num < numberOfFiles; num++ {
			wg.Add(1)

			go func(index int) {
				defer wg.Done()

				formFileForm := fmt.Sprintf("file-%d", index)
				file, err := ctx.FormFile(formFileForm)
				if err == nil {
					fileName := strings.Replace(file.Filename, " ", "_", -1)
					fileNames = append(fileNames, fileName)
					filePath := fmt.Sprintf("%s/%s", fileStoragePath, fileName)
					if err := ctx.SaveUploadedFile(file, filePath); err != nil {
						log.Printf("Error saving file %s : %v\n", fileName, err)
					}
				}
			}(num)
		}

		wg.Wait()

		// set the key
		if err := kvStore.Set(key, fileNames); err != nil {
			ctx.JSON(500, gin.H{
				"message": fmt.Sprintf("Error setting value: %v", err),
			})
			return
		}

		ctx.JSON(200, gin.H{
			"data": gin.H {
				"key": key,
			},
			"message": "File Uploaded successfully",
		})
	}
}
