package database

type Match struct {
	Title    string `json:"title"`
	Location string `json:"location"`
	HomeTeam string `json:"homeTeam"`
	AwayTeam string `json:"awayTeam"`
	Date     string `json:"date"`
}

type SimpleTeam struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Team struct {
	Name           string `json:"name"`
	Avatar         string `json:"avatar"`
	MatchesPlayed  int    `json:"matchesPlayed"`
	MatchesWon     int    `json:"matchesWon"`
	MatchesDrawn   int    `json:"matchesDrawn"`
	MatchesLost    int    `json:"matchesLost"`
	GoalsFor       int    `json:"goalsFor"`
	GoalsAgainst   int    `json:"goalsAgainst"`
	GoalDifference int    `json:"goalDifference"`
	Points         int    `json:"points"`
}

type PeacockSchedule struct {
	Name     string `json:"name"`
	HomeTeam string `json:"homeTeam"`
	AwayTeam string `json:"awayTeam"`
	Day      string `json:"day"`
	Time     string `json:"time"`
	IsReplay bool   `json:"isReplay"`
}
