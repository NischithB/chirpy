package config

import (
	"log"
	"os"

	"github.com/NischithB/chirpy/database"
	"github.com/NischithB/chirpy/models"
	"github.com/joho/godotenv"
)

var Config struct {
	DB             *database.Database[models.DatabaseModel]
	JwtSecret      []byte
	FileServerHits int
}

func ConfigureDB() {
	dbFilePath := "database.json"
	initialData := models.DatabaseModel{
		Chirps: map[int]models.Chirp{},
		Users:  map[int]models.User{},
	}
	db, err := database.NewDatabase[models.DatabaseModel](dbFilePath, initialData)
	if err != nil {
		log.Printf("Error Connecting to database: %s", err)
		os.Exit(0)
	}
	Config.DB = db
}

func LoadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Couldn't load Env variables: %s", err)
		return
	}
	Config.JwtSecret = []byte(os.Getenv("JWT_SECRET"))
}
