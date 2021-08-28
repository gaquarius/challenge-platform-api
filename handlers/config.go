package middlewares

import (
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

// DotEnvVariable -> get .env
func DotEnvVariable(key string) string {

	// load .env file
	wd, _ := os.Getwd()
	err := godotenv.Load(path.Join(wd, "/.env"))

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
