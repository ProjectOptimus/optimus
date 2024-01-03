package osc

// logging provides consistent logging objects to cut down on copy-paste across
// other packages within OpenSourceCorp

import (
	"fmt"
	"log"
	"os"
)

// Several of these vars are exported intentionally, so callers can access &
// override values during their tests or runtimes. For example, a caller
// whose functions call FatalLog() in the call stack may want to be able to
// disable IsTesting to emit logs, redirect log output to a buffer, and
// match against that output (e.g. rhad has tests like this)
var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger

	// Suppress output by default when tests are being run
	IsTesting bool
)

func init() {
	// The usage of bitwise OR here seems to be called "bitmask flagging", since
	// the log output option needs to be an integer and ORing their named bits
	// gives you a single integer result
	DebugLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	WarnLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	FatalLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
	SetLoggerPrefixName("osc")

	if os.Getenv("OSC_IS_TESTING") == "true" {
		IsTesting = true
	} else {
		IsTesting = false
	}
}

// SetLoggerPrefixName is used to construct the custom part of the log prefix of
// the loggers. It is called as part of the init() call in this package and
// given a default value, but callers can override this value with a different
// string value, typically the calling application's name (e.g. "rhad")
func SetLoggerPrefixName(name string) {
	DebugLogger.SetPrefix(fmt.Sprintf("[ %s:DEBUG ] ", name))
	InfoLogger.SetPrefix(fmt.Sprintf("[ %s:INFO  ] ", name))
	WarnLogger.SetPrefix(fmt.Sprintf("[ %s:WARN  ] ", name))
	ErrorLogger.SetPrefix(fmt.Sprintf("[ %s:ERROR ] ", name))
	FatalLogger.SetPrefix(fmt.Sprintf("[ %s:FATAL ] ", name))

}

// DebugLog throws debug log messages
func DebugLog(msg string, values ...any) {
	if !IsTesting {
		DebugLogger.Printf(msg+"\n", values...)
	}
}

// InfoLog throws info log messages
func InfoLog(msg string, values ...any) {
	if !IsTesting {
		InfoLogger.Printf(msg+"\n", values...)
	}
}

// WarnLog throws warning log messages
func WarnLog(msg string, values ...any) {
	if !IsTesting {
		WarnLogger.Printf(msg+"\n", values...)
	}
}

// ErrorLog throws error log messages
func ErrorLog(err error, msg string, values ...any) {
	if !IsTesting {
		if err != nil {
			ErrorLogger.Println(err.Error())
		}
		ErrorLogger.Printf(msg+"\n", values...)
	}
}

// FatalLog throws fatal log messages, which includes an exit call. There is no
// check here for if external tests are being run, so that callers can still see
// fatal log messages.
func FatalLog(err error, msg string, values ...any) {
	if err != nil {
		FatalLogger.Println(err.Error())
	}
	FatalLogger.Printf(msg+"\n", values...)
	if !IsTesting {
		os.Exit(1)
	}
}
