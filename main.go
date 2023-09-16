package main

import (
	"fmt"
	"log"
	"scraper/database"
	"scraper/handlers"
	"scraper/scrapers"
	"scraper/server"
)

func main() {
	fmt.Println("Starting scraper")

	scrapers.ScrapeTeams()
	scrapers.ScrapeScoreboard()
	scrapers.ScrapeMatches()

	db := database.NewDatabase()

	if err := db.LoadModels(); err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewHandler()
	handler.SetDb(db)

	server := server.NewServer()
	server.SetHandler(handler)

	server.ServerHTTP()
}
