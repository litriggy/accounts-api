package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfigs(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("unable top load .env file")
	}
	LoadApp()
	LoadDb()
}
