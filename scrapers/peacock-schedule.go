package scrapers

import (
	"encoding/json"
	"log"
	"scraper/database"
	"scraper/utils"
	"strings"

	"github.com/gocolly/colly/v2"
)

//https://www.goal.com/en-us/lists/watch-live-stream-premier-league/blt20534c45baa0e27e

func (s *Scraper) ScrapePeacockSchedule() {
	file, err := openFile("peacock-schedule.json")

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

	schedule := make([]database.PeacockSchedule, 0, 200)

	c.OnHTML("div.table-container-scroll tbody tr", func(e *colly.HTMLElement) {
		channel := e.ChildText("td:nth-child(4)")
		channels := strings.Split(channel, ",")

		if !utils.Includes("Peacock", channels) && !utils.Includes("USA Network", channels) {
			return
		}

		name := e.ChildText("td:nth-child(2)")
		teams := strings.Split(name, " vs. ")
		isReplay := utils.Includes("USA Network", channels)

		item := database.PeacockSchedule{
			Name:     name,
			HomeTeam: teams[0],
			AwayTeam: teams[1],
			Day:      e.ChildText("td:nth-child(1)"),
			Time:     e.ChildText("td:nth-child(3)"),
			IsReplay: isReplay,
		}

		schedule = append(schedule, item)

		if err := s.db.SavePeacockScheduleItem(item); err != nil {
			log.Println("Error writing peacock schedule item to database:", err)
		}
	})

	c.Visit("https://www.goal.com/en-us/lists/watch-live-stream-premier-league/blt20534c45baa0e27e")

	// Write teams to JSON file
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output
	enc.Encode(schedule)
}
