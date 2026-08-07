package main

import (
	// "github.com/Tanzor-Disco/skaM/db"
	"errors"
	"log"
	"os"

	"github.com/Tanzor-Disco/skaM/server"
	"github.com/lpernett/godotenv"
)

func getURI() (string, error) {
	options := []string{
		".env",
		"../.env",
	}
	for _, option := range options {
		err := godotenv.Load(option)
		if err == nil {
			URI := os.Getenv("URI")
			return URI, nil
		}
	}
	return "", errors.New("Didn't find a .env file")
}

func main() {
	URI, err := getURI()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.Run(URI))
}
