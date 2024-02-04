package toolkit

import "strings"

func MaskString(value string, percent float64) string {
	maskLength := int(float64(len(value)) * percent)

	if maskLength == 0 {
		return value
	}

	firstLength := (len(value) - maskLength) / 2

	return value[0:firstLength] + strings.Repeat("*", maskLength) + value[firstLength+maskLength:]
}
