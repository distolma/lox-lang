package logerror

import (
	"fmt"
	"os"
)

type LogError struct {
	HadError bool
}

func (l *LogError) Error(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
	l.HadError = true
}
