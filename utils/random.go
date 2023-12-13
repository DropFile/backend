package utils

import (
	"crypto/rand"
	"encoding/hex"
	"math"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, int(math.Ceil(float64(length) / 2)))

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	randomString := hex.EncodeToString(bytes)
	return randomString[:length], nil
}
