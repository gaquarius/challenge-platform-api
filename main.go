package main

import (
	"log"
	"net/http"

	middlewares "github.com/chattertechno/challenge-platform-api/handlers"
	"github.com/chattertechno/challenge-platform-api/routes"
	"github.com/fatih/color"
	"github.com/rs/cors"
)

func main() {
	port := middlewares.DotEnvVariable("PORT")
	color.Cyan("🌏 Server running on localhost:" + port)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	router := routes.Routes()
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	handler := c.Handler(router)
	http.ListenAndServe(":"+port, middlewares.LogRequest(handler))
}
