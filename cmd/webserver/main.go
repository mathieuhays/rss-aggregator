package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	rss "github.com/mathieuhays/rss-aggregator"
	"github.com/mathieuhays/rss-aggregator/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	config, err := rss.NewApiConfig(dbQueries)
	if err != nil {
		log.Fatal(config)
	}

	server, err := rss.NewAggregatorServer(config)
	if err != nil {
		log.Fatal(err)
	}

	_, err = rss.NewScraper(dbQueries, time.Minute)
	if err != nil {
		log.Fatal(err)
	}

	addr := "localhost:" + port

	log.Printf("Starting local server on port %s", port)
	log.Fatal(http.ListenAndServe(addr, server))
}
