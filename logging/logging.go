// Package logging provides consistent logging objects to cut down on copy-paste
// across other packages
package logging

import (
	"log"
	"os"
)

var (
	// Custom logger funcs -- actually created in init(), so they're visible everywhere else
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger

	isTesting bool
)

func init() {
	// The usage of bitwise OR here seems to be called "bitmask flagging", since
	// the log output option needs to be an integer and ORing their named bits
	// gives you a single integer result
	DebugLogger = log.New(os.Stderr, "[ rhad:DEBUG ] ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "[ rhad:INFO  ] ", log.Ldate|log.Ltime)
	WarnLogger = log.New(os.Stderr, "[ rhad:WARN  ] ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "[ rhad:ERROR ] ", log.Ldate|log.Ltime)

	if os.Getenv("RHAD_TESTING") == "true" {
		isTesting = true
	} else {
		isTesting = false
	}
}

func Debug(msg string, values ...any) {
	if !isTesting {
		DebugLogger.Printf(msg+"\n", values...)
	}
}

func Info(msg string, values ...any) {
	if !isTesting {
		InfoLogger.Printf(msg+"\n", values...)
	}
}

func Warn(msg string, values ...any) {
	if !isTesting {
		WarnLogger.Printf(msg+"\n", values...)
	}
}

func Error(msg string, values ...any) {
	if !isTesting {
		ErrorLogger.Printf(msg+"\n", values...)
	}
}
