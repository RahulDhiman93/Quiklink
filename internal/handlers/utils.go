package handlers

import (
	"github.com/skip2/go-qrcode"
	"log"
	"math/rand"
	"time"
)

// GenerateQRCode Gives back QR code given URL
func generateQRCode(url string) ([]byte, error) {
	data := url
	code, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error generating QR code: %v", err)
		return nil, err
	}

	return code, nil
}

// RandomString Generates a random string
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}
