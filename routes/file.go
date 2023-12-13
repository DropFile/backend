package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"cloud/backend/utils"

	"io"

	"github.com/gin-gonic/gin"
)

const Filestorage string = "./filestorage"
const dbURL = "http://localhost:8001"

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

		data := map[string]string{
			"key":   key,
			"value": strings.Join(fileNames, ","),
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		}

		// set the key
		resp, err := http.Post(fmt.Sprintf("%s/set", dbURL), "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error sending /set Request")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error %d - %s\n", resp.StatusCode, resp.Status)
			ctx.JSON(500, gin.H{
				"message": fmt.Sprintf("Error setting value: %v", err),
			})
			return
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

		resp, err := http.Get(fmt.Sprintf("%s/get?key=%s", dbURL, key))
		if err != nil {
			log.Println("Error sending /set Request")
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error %d - %s\n", resp.StatusCode, resp.Status)
			ctx.JSON(500, gin.H{
				"message": fmt.Sprintf("Error getting value: %v", err),
			})
			return
		}

		body, err := io.ReadAll(io.Reader(resp.Body))

		if err != nil {
			ctx.JSON(500, gin.H{
				"message": fmt.Sprintf("Error reading response: %v", err),
			})
		}

		ctx.JSON(200, gin.H{
			"data":    string(body),
			"message": "Values Found",
		})
	}
}
