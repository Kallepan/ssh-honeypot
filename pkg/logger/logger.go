package logger

import (
	"log"
	"os"
	"sync"
	"syscall"
)

var (
	CommonLog *log.Logger
	ErrorLog  *log.Logger
	FatalLog  *log.Logger
	once      sync.Once
)

func Init() {
	once.Do(func() {
		CommonLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		FatalLog = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
	})
}

// Info logs the message
func Info(msg string) {
	CommonLog.Println(msg)
}

func Infof(format string, v ...interface{}) {
	CommonLog.Printf(format, v...)
}

// Error logs the message
func Error(msg string) {
	ErrorLog.Println(msg)
}

func Errorf(format string, v ...interface{}) {
	ErrorLog.Printf(format, v...)
}

// Fatal logs the message and exits the program with exit code 1
func Fatal(msg string) {
	FatalLog.Fatalln(msg)
	syscall.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	FatalLog.Fatalf(format, v...)
	syscall.Exit(1)
}
