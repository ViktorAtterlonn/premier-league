package database

import (
	"encoding/json"
	"fmt"
	"scraper/utils"
	"sort"
	"time"
)

func (d *Database) GetScoreboard() []Team {
	scoreboard := []Team{}

	err := d.Db.Read("scoreboard", "2023-24", &scoreboard)

	if err != nil {
		fmt.Println("Error", err)
	}

	return scoreboard
}

func (d *Database) GetMatches() []Match {
	records, err := d.Db.ReadAll("matches")
	if err != nil {
		fmt.Println("Error", err)
	}

	matches := []Match{}
	for _, f := range records {
		match := Match{}
		if err := json.Unmarshal([]byte(f), &match); err != nil {
			fmt.Println("Error", err)
		}
		matches = append(matches, match)
	}

	return matches
}

func (d *Database) GetLiveMatches() []Match {
	currentTime := time.Now()
	matchDuration, _ := time.ParseDuration("90m")

	var liveMatches []Match

	for _, match := range d.matches {
		matchTime, _ := time.Parse(time.RFC3339, match.Date)

		if matchTime.Before(currentTime) && matchTime.Add(matchDuration).After(currentTime) {
			liveMatches = append(liveMatches, match)
		}
	}

	return liveMatches
}

func (d *Database) GetTeams() []SimpleTeam {
	records, err := d.Db.ReadAll("teams")
	if err != nil {
		fmt.Println("Error", err)
	}

	teams := []SimpleTeam{}
	for _, f := range records {
		team := SimpleTeam{}
		if err := json.Unmarshal([]byte(f), &team); err != nil {
			fmt.Println("Error", err)
		}
		teams = append(teams, team)
	}

	return teams
}

type UpdateTeamInput struct {
	Avatar *string `json:"avatar"`
}

func (d *Database) UpdateTeam(teamName string, input UpdateTeamInput) error {
	team := Team{}
	err := d.Db.Read("teams", teamName, &team)

	if err != nil {
		return fmt.Errorf("Team not found")
	}

	if input.Avatar != nil {
		team.Avatar = *input.Avatar
	}

	if err := d.Db.Write("teams", teamName, team); err != nil {

		fmt.Println("Error writing team to database:", err)
		return err
	}

	return nil
}

func (d *Database) SaveMatch(match Match) error {
	return d.Db.Write("matches", match.Title+match.Date, match)
}

func (d *Database) SaveTeam(team SimpleTeam) error {
	return d.Db.Write("teams", utils.GetTeamFileName(team.Name), team)
}

func (d *Database) SaveScoreboard(teams []Team) error {
	return d.Db.Write("scoreboard", "2023-24", teams)
}

func (d *Database) SavePeacockScheduleItem(schedule PeacockSchedule) error {
	return d.Db.Write("schedule", schedule.Name+schedule.Day+schedule.Time, schedule)
}

func (d *Database) GetPeacockSchedule() []PeacockSchedule {
	records, err := d.Db.ReadAll("schedule")
	if err != nil {
		fmt.Println("Error", err)
	}

	schedule := []PeacockSchedule{}
	for _, f := range records {
		item := PeacockSchedule{}
		if err := json.Unmarshal([]byte(f), &item); err != nil {
			fmt.Println("Error", err)
		}
		schedule = append(schedule, item)
	}

	return schedule
}

func (d *Database) GetAllMatches() []Match {
	records, err := d.Db.ReadAll("matches")
	if err != nil {
		fmt.Println("Error", err)
	}

	var matches []Match

	for _, f := range records {
		match := Match{}
		if err := json.Unmarshal([]byte(f), &match); err != nil {
			fmt.Println("Error", err)
		}

		if match.Date == "" {
			continue
		}

		matches = append(matches, match)
	}

	sort.Slice(matches, func(i, j int) bool {
		a, _ := time.Parse(time.RFC3339, matches[i].Date)
		b, _ := time.Parse(time.RFC3339, matches[j].Date)

		return a.Before(b)
	})

	return matches
}

func (d *Database) GetUpcomingMatches() []Match {
	currentTime := time.Now()

	var upcomingMatches []Match

	for _, match := range d.GetAllMatches() {

		matchTime, _ := time.Parse(time.RFC3339, match.Date)

		if currentTime.Before(matchTime) {
			upcomingMatches = append(upcomingMatches, match)
		}
	}

	return upcomingMatches[:8]
}

func (d *Database) GetMatchWeekNumber() string {
	matches := d.GetAllMatches()

	firstMatch := matches[0]
	lastMatch := matches[len(matches)-1]

	firstMatchTime, _ := time.Parse(time.RFC3339, firstMatch.Date)
	lastMatchTime, _ := time.Parse(time.RFC3339, lastMatch.Date)

	weekNumber := CalculateWeeks(firstMatchTime, lastMatchTime)

	return weekNumber
}

func CalculateWeeks(startDate, endDate time.Time) string {
	// Calculate the total number of weeks between startDate and endDate
	totalWeeks := int(endDate.Sub(startDate).Hours() / (24 * 7))

	// Calculate the week number for the startDate
	startWeek := int(startDate.Weekday())

	// If the startWeek is not Sunday (0), adjust it to start from Sunday
	if startWeek != 0 {
		startWeek = 7 - startWeek
	}

	// Calculate the week number
	weekNumber := (totalWeeks + startWeek) / 7

	// Magically add 3 to the week number since finals are 2 weeks long
	return fmt.Sprintf("Week %d of %d", weekNumber+1, totalWeeks+3)
}
