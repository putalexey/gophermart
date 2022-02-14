package utils

import (
	"regexp"
	"strconv"
)

var notDigitRegexp = regexp.MustCompile(`\D`)

func CheckOrderNumber(number string) bool {
	number = notDigitRegexp.ReplaceAllString(number, "")
	if number == "" {
		return false
	}

	numberLen := len(number)
	check := 0
	isOdd := numberLen%2 == 0

	for i := 0; i < numberLen; i++ {
		n, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if isOdd && i%2 == 0 {
			prod := n * 2
			if prod > 9 {
				prod -= 9
			}
			check += prod
		} else {
			check += n
		}
	}
	return check%10 == 0
}
