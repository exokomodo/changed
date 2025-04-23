package main

import (
	"fmt"
	"log"
	"os"

	"github.com/exokomodo/changed/api"
	"github.com/exokomodo/changed/pkg/changelog"
)

const (
	LocalUsername = "changed"
	LocalPassword = "changed"
	DatabaseName  = "changed"
)

func main() {
	dbURL, exists := os.LookupEnv("DATABASE_URL")
	if !exists || dbURL == "" {
		dbURL = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", LocalUsername, LocalPassword, DatabaseName)
	}

	db, err := changelog.Init(dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	repo := changelog.NewChangeRepository(db)

	log.Println("Database initialized and ready!")

	h := api.NewHandler(repo)
	router := h.SetupRouter()
	if err := api.Run(router, ":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
