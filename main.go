package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alpha2303/alpha-meetings/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	serverDomain := os.Getenv("SERVER_DOMAIN")
	serverPort := os.Getenv("SERVER_PORT")

	log.Printf("Server listening on http://%s:%s", serverDomain, serverPort)
	if err := app.StartServer(fmt.Sprintf("%s:%s", serverDomain, serverPort)); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
