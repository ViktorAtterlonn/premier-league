package scrapers

import (
	"log"
	"scraper/database"
	"scraper/utils"

	"github.com/gocolly/colly/v2"
)

func (s *Scraper) ScrapeScoreboard() {
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
			MatchesPlayed:  utils.ToInt(e.ChildText("td.widget-match-standings__matches-played")),
			MatchesWon:     utils.ToInt(e.ChildText("td.widget-match-standings__matches-won")),
			MatchesDrawn:   utils.ToInt(e.ChildText("td.widget-match-standings__matches-drawn")),
			MatchesLost:    utils.ToInt(e.ChildText("td.widget-match-standings__matches-lost")),
			GoalsFor:       utils.ToInt(e.ChildText("td.widget-match-standings__goals-for")),
			GoalsAgainst:   utils.ToInt(e.ChildText("td.widget-match-standings__goals-against")),
			GoalDifference: utils.ToInt(e.ChildText("td.widget-match-standings__goals-diff")),
			Points:         utils.ToInt(e.ChildText("td.widget-match-standings__pts")),
		}

		scoreboard = append(scoreboard, team)
	})

	c.Visit("https://www.goal.com/en/premier-league/table/2kwbbcootiqqgmrzs6o5inle5")

	if err := s.db.SaveScoreboard(scoreboard); err != nil {
		log.Println("Error writing scoreboard to database:", err)
	}
}
