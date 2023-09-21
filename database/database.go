package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	scribble "github.com/nanobox-io/golang-scribble"
)

type Database struct {
	Db         *scribble.Driver
	matches    []Match
	teams      []Team
	scoreboard []Team
	schedule   []PeacockSchedule
}

func NewDatabase() *Database {
	db, err := scribble.New("./local", nil)

	if err != nil {
		log.Panicf("Error creating database: %s", err)
	}

	return &Database{
		Db:         db,
		matches:    []Match{},
		teams:      []Team{},
		scoreboard: []Team{},
		schedule:   []PeacockSchedule{},
	}
}

func (d *Database) LoadModels() error {
	var matches []Match
	var scoreboard []Team
	var teams []Team
	var schedule []PeacockSchedule

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

	scheduleData, err := readFile("peacock-schedule.json")
	if err != nil {
		fmt.Println("Error reading peacock-schedule.json:", err)
		return err
	}

	if err := json.Unmarshal(scheduleData, &schedule); err != nil {
		fmt.Println("Error unmarshalling peacock-schedule.json:", err)
		return err
	}

	d.matches = matches
	d.teams = teams
	d.scoreboard = scoreboard
	d.schedule = schedule

	fmt.Println("✨ Database loaded")

	fmt.Printf("Matches: %d\n", len(matches))
	fmt.Printf("Teams: %d\n", len(teams))
	fmt.Printf("Scoreboard: %d\n", len(scoreboard))
	fmt.Printf("Peacock Schedule: %d\n", len(schedule))

	return nil
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
