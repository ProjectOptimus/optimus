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
)

func init() {
	// The usage of bitwise OR here seems to be called "bitmask flagging", since
	// the log output option needs to be an integer and ORing their named bits
	// gives you a single integer result
	DebugLogger = log.New(os.Stderr, "[ DEBUG ] ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "[ INFO  ] ", log.Ldate|log.Ltime)
	WarnLogger = log.New(os.Stderr, "[ WARN  ] ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "[ ERROR ] ", log.Ldate|log.Ltime)
}

func Debug(msg string, values ...any) {
	DebugLogger.Printf(msg+"\n", values...)
}

func Info(msg string, values ...any) {
	InfoLogger.Printf(msg+"\n", values...)
}

func Warn(msg string, values ...any) {
	WarnLogger.Printf(msg+"\n", values...)
}

func Error(msg string, values ...any) {
	ErrorLogger.Printf(msg+"\n", values...)
}
