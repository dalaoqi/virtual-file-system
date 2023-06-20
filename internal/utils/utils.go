package utils

import "regexp"

var (
	re *regexp.Regexp
)

func init() {
	// Regular expression to match invalid characters
	re = regexp.MustCompile(`(?m)[^a-zA-Z0-9]+`)
}

// Check if the given string contains invalid characters
func ExistInvalidChars(str string) bool {
	match := re.FindAllString(str, -1)
	return len(match) > 0
}
