package helpers

func SerialNumbersEndsWithOddDigit(serialNumber string) bool {
	if len(serialNumber) == 0 {
		return false
	}

	lastChar := serialNumber[len(serialNumber)-1]
	return lastChar%2 != 0
}
