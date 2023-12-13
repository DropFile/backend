package utils

import "os"

func CreateFolder (folderPath string) error {
	err := os.Mkdir(folderPath, 0755)
	return err
}