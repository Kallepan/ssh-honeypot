package conf

import (
	"log"
	"os"
	"path/filepath"
)

var (
	// Log is the logger
	Log      *log.Logger
	ErrorLog *log.Logger
)

func StartLogger() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filepath.Join(wd, "logs", "ssh.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)
}
