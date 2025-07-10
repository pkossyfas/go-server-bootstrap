/*
Package logger is a json structured log provider.
*/
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LogMessage is the basic structure used for log messages.
type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message,omitempty"`
	Exception string    `json:"exception,omitempty"`
}

var lg = &LogMessage{}

// Log structures the log message to a json format.
func Log(message string) {
	defer func() {
		lg.Exception = ""
		lg.Message = ""
	}()
	lg.Timestamp = time.Now()
	lg.Message = message
	l, err := json.Marshal(lg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(l))
}

// Info logs a string info message.
func Info(format string, a ...interface{}) {
	lg.Level = "info"
	Log(fmt.Sprintf(format, a...))
}

// Warn logs a warn message and the respective exception.
func Warn(err error, format string, a ...interface{}) {
	lg.Level = "warn"
	if err != nil {
		lg.Exception = err.Error()
	} else {
		lg.Exception = ""
	}
	Log(fmt.Sprintf(format, a...))
}

// Error logs an error message and the respective exception.
func Error(err error, format string, a ...interface{}) {
	lg.Level = "error"
	if err != nil {
		lg.Exception = err.Error()
	} else {
		lg.Exception = ""
	}
	Log(fmt.Sprintf(format, a...))
}

// Fatal logs a fatal exception and then exits.
func Fatal(err error) {
	lg.Level = "fatal"
	lg.Exception = err.Error()
	Log("")
	os.Exit(1)
}
