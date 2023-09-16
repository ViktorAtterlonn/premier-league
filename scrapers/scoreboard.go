package scrapers

import (
	"encoding/json"
	"log"
	"scraper/database"
	"strconv"

	"github.com/gocolly/colly/v2"
)

func ScrapeScoreboard() {
	//https://www.premierleague.com/fixtures

	file, err := openFile("scoreboard.json")

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

	scoreboard := make([]database.Team, 0, 200)

	c.OnHTML("tr.widget-match-standings__row", func(e *colly.HTMLElement) {
		log.Println(e.ChildText("span.widget-match-standings__team--full-name"))

		name := e.ChildText("span.widget-match-standings__team--full-name")

		team := database.Team{
			Name:           name,
			Avatar:         "/public/images/teams/" + getTeamFileName(name) + getTeamFileExtension((name)),
			MatchesPlayed:  toInt(e.ChildText("td.widget-match-standings__matches-played")),
			MatchesWon:     toInt(e.ChildText("td.widget-match-standings__matches-won")),
			MatchesDrawn:   toInt(e.ChildText("td.widget-match-standings__matches-drawn")),
			MatchesLost:    toInt(e.ChildText("td.widget-match-standings__matches-lost")),
			GoalsFor:       toInt(e.ChildText("td.widget-match-standings__goals-for")),
			GoalsAgainst:   toInt(e.ChildText("td.widget-match-standings__goals-against")),
			GoalDifference: toInt(e.ChildText("td.widget-match-standings__goals-diff")),
			Points:         toInt(e.ChildText("td.widget-match-standings__pts")),
		}

		scoreboard = append(scoreboard, team)

	})

	c.Visit("https://www.goal.com/en/premier-league/table/2kwbbcootiqqgmrzs6o5inle5")

	// Write teams to JSON file
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(scoreboard)
}

func toInt(str string) int {
	intVar, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}

	return intVar
}
