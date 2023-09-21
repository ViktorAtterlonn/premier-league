package jobs

import (
	"log"
)

func (r *Runner) JobCheckScoreboard() {
	log.Println("Job started: CheckScoreboard")
	r.scraper.ScrapeScoreboard()
	log.Println("Job finished: CheckScoreboard")
}
