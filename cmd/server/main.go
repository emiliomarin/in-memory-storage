package main

import (
	"log"
	"os"

	"in-memory-storage/internal/app"
)

func main() {
	application, err := app.New(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatal("error creating application: ", err)
	}

	application.Start()
}
