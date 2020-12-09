package security

import (
	"math"
	"strings"
)

const (
	// Base64Characters list of base64 characters
	Base64Characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
)

// ShannonEntropy calculate shannon entropy
func ShannonEntropy(data string) (entropy float64) {
	if data == "" {
		return 0
	}
	charCounts := make(map[rune]int)
	for _, char := range data {
		charCounts[char]++
	}
	invLength := 1.0 / float64(len(data))
	for _, count := range charCounts {
		freq := float64(count) * invLength
		entropy -= freq * math.Log2(freq)
	}
	return entropy
}

// ShannonEntropyBase64 calculate shannon entropy
func ShannonEntropyBase64(data string) (entropy float64) {
	if data == "" {
		return 0
	}
	for _, character := range Base64Characters {
		probability := float64(strings.Count(data, string(character)) / len(data))
		if probability > 0 {
			entropy += -probability * math.Log2(probability)
		}
	}
	return entropy
}
