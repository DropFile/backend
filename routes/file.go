package routes

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"cloud/backend/utils"

	"github.com/gin-gonic/gin"
)

const Filestorage string = "./filestorage"
const jsonFilePath string = "./db.json"

func HandleUpload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, err := utils.GenerateRandomString(20)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Error generating Random String",
			})
			return
		}
		fmt.Println(key)

		fileStoragePath := fmt.Sprintf("%s/%s", Filestorage, key)
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
					fileExtension, fileError := utils.GetFileExtension(fileName)

					if fileError != nil {
						log.Printf("%s: %v\n", fileName, fileError)
					}

					if fileExtension == "mp4" {
						fileNames = append(fileNames, fileName)
						filePath := fmt.Sprintf("%s/%s", fileStoragePath, fileName)
						err := ctx.SaveUploadedFile(file, filePath)
						if err != nil {
							log.Printf("Error saving file %s : %v\n", fileName, err)
						} else {
							err := utils.SegmentVideoFile(filePath, fileName)
							if err != nil {
								log.Printf("Error Segmenting File %s : %v\n", fileName, err)
							}
						}
					}

				}
			}(num)
		}

		wg.Wait()

		combinedValue := strings.Join(fileNames, ",")
		fmt.Println(combinedValue)
		err = utils.WriteJson(jsonFilePath, key, combinedValue)
		if err != nil {
			log.Fatal("Error writing to JSON")
		}

		ctx.JSON(200, gin.H{
			"data": gin.H{
				"key": key,
			},
			"message": "File Uploaded successfully",
		})
	}
}

func HandleFileMetadata() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Query("key")

		result, err := utils.GetValueForKey(jsonFilePath, key)

		if err != nil {
			ctx.JSON(500, gin.H{
				"message": "Can't find any value",
			})
		}

		ctx.JSON(200, gin.H{
			"data":    string(result),
			"message": "Values Found",
		})
	}
}
