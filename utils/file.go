package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func CreateFolder(folderPath string) error {
	err := os.Mkdir(folderPath, 0755)
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
		"-b:a", "128k",
		"-map", "0:0",
		"-f", "segment",
		"-segment_time", "10",
		"-segment_list", fmt.Sprintf("%slist.m3u8", filePath),
		"-segment_format", "mp4",
		fmt.Sprintf("%s%%03d.%s", filePath, "mp4"),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
