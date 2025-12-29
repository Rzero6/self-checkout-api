package utils

import (
	"strconv"
	"strings"
)

func FormatPrice(amount int64) string {
	s := strconv.FormatInt(amount, 10)

	var result []string
	for len(s) > 3 {
		result = append([]string{s[len(s)-3:]}, result...)
		s = s[:len(s)-3]
	}
	result = append([]string{s}, result...)

	return "Rp " + strings.Join(result, ".")
}
