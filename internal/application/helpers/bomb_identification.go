package helpers

import (
	"strings"
	"unicode"
)

func SerialNumbersEndsWithOddDigit(serialNumber string) bool {
	if len(serialNumber) == 0 {
		return false
	}

	lastChar := serialNumber[len(serialNumber)-1]
	return lastChar%2 != 0
}

func SerialNumberContainsVowel(serialNumber string) bool {
	vowels := "aeiou"
	for _, char := range serialNumber {
		if strings.ContainsRune(vowels, unicode.ToLower(char)) {
			return true
		}
	}
	return false
}
