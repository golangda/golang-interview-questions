package word

import "strings"

// Count counts words (split by whitespace).
func Count(s string) int {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	return len(strings.Fields(s))
}

// RepeatCount repeats Count several times (simulate heavy work)
func RepeatCount(s string, times int) int {
	total := 0
	for i := 0; i < times; i++ {
		total += Count(s)
	}
	return total
}
