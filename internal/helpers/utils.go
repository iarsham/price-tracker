package helpers

import "strings"

func ThousandSeparator(s string) string {
	parts := strings.Split(s, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 {
		decimalPart = "." + parts[1]
	}
	n := len(integerPart)
	for i := n - 3; i > 0; i -= 3 {
		integerPart = integerPart[:i] + "," + integerPart[i:]
	}
	return integerPart + decimalPart
}
