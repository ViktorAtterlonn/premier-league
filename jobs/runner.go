package jobs

import (
	"scraper/scrapers"

	"github.com/robfig/cron"
)

type Runner struct {
	scraper *scrapers.Scraper
	c       *cron.Cron
}

func NewRunner() *Runner {
	return &Runner{
		c: cron.New(),
	}
}

func (r *Runner) SetScraper(scraper *scrapers.Scraper) {
	r.scraper = scraper
}

func (r *Runner) SetupJobs() {
	r.c.AddFunc("0 30 * * * *", r.JobCheckScoreboard)
}

func (r *Runner) Run() {
	r.c.Start()
}
