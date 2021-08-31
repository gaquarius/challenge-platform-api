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
	log.Println(path.Join(wd, "/.env"))

	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
