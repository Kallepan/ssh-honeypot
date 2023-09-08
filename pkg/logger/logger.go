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
		logFile, err := os.OpenFile(os.Getenv("LOGS_FILE"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %s", err)
		}

		CommonLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLog = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		FatalLog = log.New(logFile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
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
