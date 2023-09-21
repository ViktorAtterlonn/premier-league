package scrapers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"scraper/database"
	"strings"

	"github.com/gocolly/colly/v2"
)

func (s *Scraper) ScrapeTeams() {

	file, err := openFile("teams.json")

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

	teams := make([]database.SimpleTeam, 0, 200)

	c.OnHTML("li.clubList__club", func(e *colly.HTMLElement) {

		team := database.SimpleTeam{
			Name:   e.ChildText("span.name"),
			Avatar: e.ChildAttr("img.badge-image.badge-image--50.js-badge-image", "src"),
		}

		// Download imamge
		downloadImage(team.Avatar, getTeamFileName(team.Name))

		teams = append(teams, team)

		if err := s.db.SaveTeam(team); err != nil {
			fmt.Println("Error writing team to database:", err)
		}
	})

	c.Visit("https://www.premierleague.com/fixtures")
}

func getTeamFileName(teamName string) string {
	return strings.Replace(strings.ToLower(teamName), " ", "-", -1)
}

func getTeamFileExtension(teamName string) string {
	teamFileName := getTeamFileName(teamName)
	exception := []string{"manchester-united", "newcastle-united"}

	if contains(exception, getTeamFileName(teamFileName)) {
		return ".png"
	}

	return ".svg"
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func downloadImage(url string, fileName string) error {
	// Extract the filename from the URL
	tokens := strings.Split(url, "/")
	extenstion := strings.Split(tokens[len(tokens)-1], ".")[1]

	// Create or overwrite the file
	file, err := os.Create("templates/public/images/teams/" + fileName + "." + extenstion)
	if err != nil {
		return err
	}
	defer file.Close()

	// Send an HTTP GET request to the image URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the image content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Image downloaded:", fileName)
	return nil
}
