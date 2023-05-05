package logger

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
)

var (
	// Log is the logger
	logger      *log.Logger
	errorLogger *log.Logger
)

func Logf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func Log(format string) {
	logger.Printf(format)
}

func Fatal(message string) {
	errorLogger.Fatal(message)
	syscall.Exit(1)
}

func Errf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

func Err(message string) {
	errorLogger.Println(message)
}

func StartLogger() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filepath.Join(wd, "logs.txt"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(file, "", log.Ldate|log.Ltime)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
