package scrapers

import (
	"fmt"
	"log"
	"net/url"
	"scraper/database"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func (s *Scraper) ScrapeMatches() {
	//https://www.premierleague.com/fixtures

	file, err := openFile("matches.json")

	if err != nil {
		return
	}

	defer file.Close()

	// Instantiate default collector
	c := getCachedCollector()

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	matches := make([]database.Match, 0, 200)

	c.OnHTML("div.match-row", func(e *colly.HTMLElement) {

		match := database.Match{
			Title:    e.ChildAttr("meta[itemprop=name]", "content"),
			Location: e.ChildAttr("meta[itemprop=location]", "content"),
			HomeTeam: e.ChildText("div.team-home > span.team-name"),
			AwayTeam: e.ChildText("div.team-away > span.team-name"),
			Date:     parseDateToCETDate(e.ChildAttr("meta[itemprop=startDate]", "content")),
		}

		matches = append(matches, match)

		if err := s.db.SaveMatch(match); err != nil {
			fmt.Println("Error writing match to database:", err)
		}
	})

	visited := make(map[string]bool)

	c.OnHTML("a.nav-switch__next", func(e *colly.HTMLElement) {
		link := e.Attr("data-href")
		nextUrl := getNextFixuresDay(link)

		if visited[nextUrl] {
			return
		}

		// if len(visited) > 20 {
		// 	return
		// }

		visited[nextUrl] = true

		log.Println("Next page found: ", getNextFixuresDay(link))

		e.Request.Visit(nextUrl)
	})

	c.Visit("https://www.goal.com/en/premier-league/fixtures-results/2kwbbcootiqqgmrzs6o5inle5")
}

// 2023-09-23T14:00:00Z -> 2023-09-23 16:00:00
func parseDateToCETDate(timeValue string) string {
	etLocation, err := time.LoadLocation("America/New_York")

	if err != nil {
		fmt.Println("Error loading Eastern Time location:", err)
		return ""
	}

	etTime, err := time.Parse(time.RFC3339, timeValue)

	if err != nil {
		fmt.Println("Error parsing time:", err)
		return ""
	}

	etTime = etTime.In(etLocation)

	// Convert the time to Central European Time (CET)
	cetLocation, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Println("Error loading Central European Time location:", err)
		return ""
	}
	cetTime := etTime.In(cetLocation)

	return cetTime.Format("2006-01-02T15:04:05Z")
}

// https://www.goal.com/en/ajax/competition/matches/2kwbbcootiqqgmrzs6o5inle5?date=2023-09-18%2018%3A45%3A00&calendarId=1jt5mxgn4q5r6mknmlqv5qjh0
func getNextFixuresDay(href string) string {
	searchParams, err := url.Parse(href)
	if err != nil {
		return ""
	}

	values := searchParams.Query()

	date := strings.Split(values.Get("date"), " ")[0]

	const pattern = "https://www.goal.com/en/premier-league/fixtures-results/{date}/2kwbbcootiqqgmrzs6o5inle5"

	return strings.Replace(pattern, "{date}", date, 1)
}
