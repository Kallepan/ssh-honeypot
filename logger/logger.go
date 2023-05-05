package logger

import (
	"log"
	"os"
	"path/filepath"
)

func Start() {
	file, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), "ssh", ".log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Println("Starting SSH Logger")
}
