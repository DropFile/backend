package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, int(math.Ceil(float64(length) * 3 / 4)))

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	randomString := base64.URLEncoding.EncodeToString(bytes)
	return randomString[:length], nil
}
