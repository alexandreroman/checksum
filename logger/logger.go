// Simple application logger.
package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	lineTerminator = "\n"
)

// Internal lock used for logging.
var lock sync.Mutex = sync.Mutex{}

// Set to true to enable logs.
var Verbose bool = false

// Log an entry.
// This entry is only displayed in verbose mode.
func Debug(format string, a ...interface{}) {
	if Verbose {
		print(os.Stderr, format, a...)
	}
}

// Log an entry.
func Info(format string, a ...interface{}) {
	print(os.Stdout, format, a...)
}

// Log an entry, and exit this program.
func Fatal(format string, a ...interface{}) {
	buf := new(bytes.Buffer)
	print(buf, format, a...)
	panic(buf)
}

func print(out io.Writer, format string, a ...interface{}) {
	// Add a new line terminator if needed.
	if !strings.HasSuffix(format, lineTerminator) {
		format += lineTerminator
	}

	// Acquire a lock to make sure entries are not mixed up.
	lock.Lock()
	// The lock will automatically be released.
	defer lock.Unlock()

	// Do the logging.
	fmt.Fprintf(out, format, a...)
}
