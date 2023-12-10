package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/viquitorreis/skills-back.git/auth"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error fetching port in .env")
	}

	auth.NewAuth()
	server := NewApiServer(port, store)
	server.Run()
}
