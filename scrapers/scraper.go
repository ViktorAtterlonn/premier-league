package scrapers

import (
	"log"
	"os"
	"scraper/database"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	db *database.Database
}

func NewScraper() *Scraper {
	return &Scraper{}
}

func (s *Scraper) SetDb(db *database.Database) {
	s.db = db
}

func openFile(fName string) (*os.File, error) {
	file, err := os.Create("database/" + fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return nil, err
	}
	return file, nil
}

func getCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("goal.com", "www.goal.com", "www.premierleague.com", "premierleague.com"),
	)

	return c
}

func getCachedCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("goal.com", "www.goal.com", "www.premierleague.com", "premierleague.com"),
		colly.CacheDir("./premierleague_cache"),
	)

	return c
}
