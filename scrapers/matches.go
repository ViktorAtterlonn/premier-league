package scrapers

import (
	"encoding/json"
	"log"
	"net/url"
	"scraper/database"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ScrapeMatches() {
	//https://www.premierleague.com/fixtures

	file, err := openFile("matches.json")

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

	matches := make([]database.Match, 0, 200)

	c.OnHTML("div.match-row", func(e *colly.HTMLElement) {

		match := database.Match{
			Title:    e.ChildAttr("meta[itemprop=name]", "content"),
			Location: e.ChildAttr("meta[itemprop=location]", "content"),
			HomeTeam: e.ChildText("div.team-home > span.team-name"),
			AwayTeam: e.ChildText("div.team-away > span.team-name"),
			Date:     e.ChildAttr("meta[itemprop=startDate]", "content"),
		}

		matches = append(matches, match)
	})

	visited := make(map[string]bool)

	c.OnHTML("a.nav-switch__next", func(e *colly.HTMLElement) {
		link := e.Attr("data-href")
		nextUrl := getNextFixuresDay(link)

		if visited[nextUrl] {
			return
		}

		if len(visited) > 20 {
			return
		}

		visited[nextUrl] = true

		log.Println("Next page found: ", getNextFixuresDay(link))

		e.Request.Visit(nextUrl)
	})

	c.Visit("https://www.goal.com/en/premier-league/fixtures-results/2kwbbcootiqqgmrzs6o5inle5")

	// Write teams to JSON file
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(matches)
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
