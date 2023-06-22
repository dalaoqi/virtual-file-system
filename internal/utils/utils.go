package utils

import (
	"regexp"
	"strings"
)

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

// Split the input into individual arguments while preserving quoted strings
func SplitArguments(line string) []string {
	var args []string
	fields := strings.Fields(line)
	quoteOpen := false
	var currentArg string

	for _, field := range fields {
		if strings.HasPrefix(field, "\"") && strings.HasSuffix(field, "\"") {
			args = append(args, strings.TrimSuffix(strings.TrimPrefix(field, "\""), "\""))
		} else if strings.HasPrefix(field, "\"") && !quoteOpen {
			quoteOpen = true
			currentArg = strings.TrimPrefix(field, "\"")
		} else if strings.HasSuffix(field, "\"") && quoteOpen {
			quoteOpen = false
			currentArg += " " + strings.TrimSuffix(field, "\"")
			args = append(args, currentArg)
		} else if quoteOpen {
			currentArg += " " + field
		} else {
			args = append(args, field)
		}
	}
	return args
}
