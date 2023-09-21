package main

import (
	"fmt"
	"scraper/database"
	"scraper/handlers"
	"scraper/jobs"
	"scraper/scrapers"
	"scraper/server"
)

func main() {
	fmt.Println("Starting scraper")

	db := database.NewDatabase()

	scraper := scrapers.NewScraper()
	scraper.SetDb(db)

	scraper.ScrapeTeams()
	scraper.ScrapeScoreboard()
	scraper.ScrapeMatches()
	scraper.ScrapePeacockSchedule()

	runner := jobs.NewRunner()
	runner.SetScraper(scraper)
	runner.SetupJobs()
	runner.Run()

	handler := handlers.NewHandler()
	handler.SetDb(db)

	server := server.NewServer()
	server.SetHandler(handler)

	server.ServerHTTP()
}
