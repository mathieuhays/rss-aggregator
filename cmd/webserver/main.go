package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	rss "github.com/mathieuhays/rss-aggregator"
	"log"
	"net/http"
	"os"
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

	config, err := rss.NewApiConfig(db)
	if err != nil {
		log.Fatal(config)
	}

	server, err := rss.NewAggregatorServer(config)
	if err != nil {
		log.Fatal(err)
	}

	addr := "localhost:" + port

	log.Printf("Starting local server on port %s", port)
	log.Fatal(http.ListenAndServe(addr, server))
}
