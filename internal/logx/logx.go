package logx

import (
	"fmt"
	"os"
)

var isVerbose bool

// EnableVerboseOutput enables debug logging
func EnableVerboseOutput() {
	isVerbose = true
}

// Info prints an info to stderr
// Most message should be log at this level
func Info(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

// Debug prints an debug to stderr in verbose mode
func Debug(format string, v ...interface{}) {
	if isVerbose {
		fmt.Fprintf(os.Stderr, format+"\n", v...)
	}
}
