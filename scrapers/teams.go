package scrapers

import (
	"encoding/json"
	"log"
	"scraper/database"

	"github.com/gocolly/colly/v2"
)

func ScrapeTeams() {

	file, err := openFile("teams.json")

	if err != nil {
		return
	}

	defer file.Close()

	// Instantiate default collector
	c := getCollector()

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	teams := make([]database.SimpleTeam, 0, 200)

	c.OnHTML("li.clubList__club", func(e *colly.HTMLElement) {

		team := database.SimpleTeam{
			Name:   e.ChildText("span.name"),
			Avatar: e.ChildAttr("img.badge-image.badge-image--50.js-badge-image", "src"),
		}

		teams = append(teams, team)
	})

	c.Visit("https://www.premierleague.com/fixtures")

	// Write teams to JSON file
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(teams)
}
