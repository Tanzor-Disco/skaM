package main

import (
	// "github.com/Tanzor-Disco/skaM/db"
	"log"
	"os"

	"github.com/Tanzor-Disco/skaM/db/migrations"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/Tanzor-Disco/skaM/server"

	"github.com/lpernett/godotenv"
)

func main() {
	//Load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	//Load .env variables
	URI := os.Getenv("URI")
	baseURL := os.Getenv("BASE_URL")
	SMTPUsername := os.Getenv("SMTP_USERNAME")
	SMTPPassword := os.Getenv("SMTP_PASSWORD")
	SMTPHost := os.Getenv("SMTP_HOST")
	SMTPAddr := os.Getenv("SMTP_ADDRESS")
	SMTPFrom := os.Getenv("SMTP_FROM")

	//Update migrations
	err := migrations.Update(URI)
	if err != nil {
		log.Fatal(err)
	}

	//Run the server
	SMTPData := models.NewSMTPData(SMTPUsername, SMTPPassword, SMTPHost, SMTPAddr, SMTPFrom, baseURL)
	serverData := models.NewServerData(URI, SMTPData)
	log.Fatal(server.Run(serverData))
}
