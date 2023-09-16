package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Database struct {
	matches    []Match
	teams      []Team
	scoreboard []Team
}

func NewDatabase() *Database {
	return &Database{
		matches:    []Match{},
		teams:      []Team{},
		scoreboard: []Team{},
	}
}

func (d *Database) LoadModels() error {
	var matches []Match
	var scoreboard []Team
	var teams []Team

	// Load models
	matchesData, err := readFile("matches.json")
	if err != nil {
		fmt.Println("Error reading matches.json:", err)
		return err
	}

	// Unmarshal the JSON data into the map
	if err := json.Unmarshal(matchesData, &matches); err != nil {
		fmt.Println("Error unmarshalling matches.json:", err)
		return err
	}

	teamsData, err := readFile("teams.json")
	if err != nil {
		fmt.Println("Error reading teams.json:", err)
	}

	// Unmarshal the JSON data into the map
	if err := json.Unmarshal(teamsData, &teams); err != nil {
		fmt.Println("Error unmarshalling teams.json:", err)
		return err
	}

	scoreboardData, err := readFile("scoreboard.json")
	if err != nil {
		fmt.Println("Error reading scoreboard.json:", err)
		return err
	}

	if err := json.Unmarshal(scoreboardData, &scoreboard); err != nil {
		fmt.Println("Error unmarshalling teams.json:", err)
		return err
	}

	d.matches = matches
	d.teams = teams
	d.scoreboard = scoreboard

	fmt.Println("✨ Database loaded")

	fmt.Printf("Matches: %d\n", len(matches))
	fmt.Printf("Teams: %d\n", len(teams))
	fmt.Printf("Scoreboard: %d\n", len(scoreboard))

	return nil
}

func (d *Database) GetMatches() []Match {
	return d.matches
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

func (d *Database) GetTeams() []Team {
	return d.teams
}

func (d *Database) GetScoreboard() []Team {
	return d.scoreboard
}

func readFile(name string) ([]byte, error) {
	file, err := os.Open("database/" + name)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	defer file.Close()

	// Read the file content into a []byte
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// Create a map to hold the parsed JSON data

	return data, nil
}
