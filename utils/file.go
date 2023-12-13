package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func CreateFolder(folderPath string) error {
	fmt.Println(folderPath)
	err := os.MkdirAll(folderPath, 0755)
	return err
}

func GetFileExtension(filename string) (string, error) {
	index := len(filename) - 1
	for ; index >= 0; index-- {
		if filename[index] == '.' {
			break
		}
	}

	if index == -1 {
		return "", errors.New("invalid file format")
	}

	fileExtension := ""
	for num := index + 1; num < len(filename); num++ {
		fileExtension = fileExtension + string(filename[num])
	}

	return fileExtension, nil
}

func GetFileWithoutExtension(filename string) (string, error) {
	index := len(filename) - 1
	for ; index >= 0; index-- {
		if filename[index] == '.' {
			break
		}
	}

	if index == -1 {
		return "", errors.New("invalid file format")
	}

	fileNameWithoutExtension := ""
	for num := 0; num < index; num++ {
		fileNameWithoutExtension = fileNameWithoutExtension + string(filename[num])
	}

	fmt.Println(fileNameWithoutExtension)
	return fileNameWithoutExtension, nil
}

func SegmentVideoFile(filePath string, fileName string) error {
	cmd := exec.Command("ffmpeg",
		"-i", filePath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-profile:v", "baseline",
		"-b:v", "400k",
		"-b:a", "64k",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-s", "640x360",
		"-start_number", "0",
		filePath+".m3u8",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		cmd := exec.Command("ffmpeg",
			"-i", filePath+".m3u8",
			"-c", "copy",
			"-bsf:a", "aac_adtstoasc",
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-start_number", "0",
			"-f", "hls",
			filePath+"playlist.m3u8",
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
	}
	return err
}

func WriteJson(jsonFilePath, key, value string) error {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatal("Error opening JSON file:", err)
	}
	defer file.Close()
	jsonReader := json.NewDecoder(file)
	var data map[string]interface{}
	if err := jsonReader.Decode(&data); err != nil {
		log.Fatal("Error decoding JSON:", err)
	}
	data[key] = value
	writeFile, err := os.Create(jsonFilePath)
	if err != nil {
		log.Fatal("Error opening JSON file for writing:", err)
	}
	defer writeFile.Close()
	jsonWriter := json.NewEncoder(writeFile)
	if err := jsonWriter.Encode(data); err != nil {
		log.Fatal("Error encoding and writing JSON:", err)
	}
	return nil
}

func GetValueForKey(jsonFilePath string, key string) (string, error) {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()
	jsonReader := json.NewDecoder(file)
	var data map[string]interface{}
	if err := jsonReader.Decode(&data); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}
	value, exists := data[key]
	if !exists {
		return "", fmt.Errorf("key %s not found in JSON", key)
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("value for key %s is not a string", key)
	}
	return strValue, nil
}
