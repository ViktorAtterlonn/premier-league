package handlers

import (
	"fmt"
	"net/http"
	"scraper/database"
	"scraper/utils"
	"strings"
	"text/template"
	"time"

	"github.com/unrolled/render"
)

type Handler struct {
	db *database.Database
	r  *render.Render
}

func NewHandler() *Handler {
	return &Handler{
		r: render.New(render.Options{
			Extensions: []string{".html"},
			Layout:     "layout",
			Funcs: []template.FuncMap{
				{
					"getMatchTime":             getMatchTime,
					"getTeamImage":             getTeamImage,
					"formatDate":               formatDate,
					"getMatchElapsedPercent":   getMatchElapsedPercent,
					"getMatchRemainingPercent": getMatchRemainingPercent,
					"getDoesAirOnPeacock":      getDoesAirOnPeacock,
					"getIsReplayOnPeacock":     getIsReplayOnPeacock,
					"getTeamPosition":          getTeamPosition,
				},
			},
		}),
	}
}

func formatDate(timestamp string) string {
	inputTime, _ := time.Parse(time.RFC3339, timestamp)
	// Format the date as "Saturday 23/9 16:00"
	formattedDate := inputTime.Format("Mon 2/1 15:04")
	return formattedDate
}

func getTeamImage(name string) string {
	return "/public/images/teams/" + utils.GetTeamFileName(name) + utils.GetTeamFileExtension((name))
}

// timestamp is a string in the format of "2021-09-19T15:00:00Z"
func getMatchTime(timestamp string) string {
	// Parse timestamp into a time.Time object
	startDate, _ := time.Parse(time.RFC3339, timestamp)
	currentTime := time.Now()

	if currentTime.Before(startDate) {
		return "Not started"
	}

	if startDate.Add(time.Hour * 2).Before(currentTime) {
		return "Finished"
	}

	// Calculate the elapsed time since the start date
	elapsedTime := currentTime.Sub(startDate)
	minutes := int(elapsedTime.Minutes())
	half := "1st half"

	// Check if the game is in the 2nd half
	if minutes >= 45 {
		half = "2nd half"
	}

	// Calculate the remaining minutes in the current half
	remainingMinutes := minutes % 45

	// Format the result
	if remainingMinutes > 0 {
		return fmt.Sprintf("%s %d'", half, remainingMinutes)
	}

	return fmt.Sprintf("%s %d", half, minutes)
}

func getMatchRemainingPercent(timestamp string) int {
	// Parse timestamp into a time.Time object
	startDate, _ := time.Parse(time.RFC3339, timestamp)
	currentTime := time.Now()

	if currentTime.Before(startDate) {
		return 0
	}

	if startDate.Add(time.Hour * 2).Before(currentTime) {
		return 100
	}

	// Calculate the elapsed time since the start date
	elapsedTime := currentTime.Sub(startDate)
	minutes := int(elapsedTime.Minutes())

	return minutes / 90 * 100
}

func getMatchElapsedPercent(timestamp string) int {
	// Parse timestamp into a time.Time object
	startDate, _ := time.Parse(time.RFC3339, timestamp)
	currentTime := time.Now()

	if currentTime.Before(startDate) {
		return 0
	}

	if startDate.Add(time.Hour * 2).Before(currentTime) {
		return 100
	}

	// Calculate the elapsed time since the start date
	elapsedTime := currentTime.Sub(startDate)
	minutes := int(elapsedTime.Minutes())

	return minutes / 90 * 100
}

func getDoesAirOnPeacock(match database.Match, schedule []database.PeacockSchedule) bool {
	for _, item := range schedule {
		isHomeTeam := strings.Contains(match.HomeTeam, item.HomeTeam) || strings.Contains(getTeamNickName(match.HomeTeam), item.HomeTeam)
		isAwayTeam := strings.Contains(match.AwayTeam, item.AwayTeam) || strings.Contains(getTeamNickName(match.AwayTeam), item.AwayTeam)
		isSameTime := match.Date == item.Date

		if isHomeTeam && isAwayTeam && isSameTime {
			return true
		}
	}

	return false
}

func getIsReplayOnPeacock(match database.Match, schedule []database.PeacockSchedule) bool {
	for _, item := range schedule {
		isHomeTeam := strings.Contains(match.HomeTeam, item.HomeTeam) || strings.Contains(getTeamNickName(match.HomeTeam), item.HomeTeam)
		isAwayTeam := strings.Contains(match.AwayTeam, item.AwayTeam) || strings.Contains(getTeamNickName(match.AwayTeam), item.AwayTeam)
		isSameTime := match.Date == item.Date

		if isHomeTeam && isAwayTeam && isSameTime {
			return item.IsReplay
		}
	}

	return false
}

func getTeamNickName(teamName string) string {
	if teamName == "Wolverhampton Wanderers" {
		return "Wolves"
	}

	return teamName
}

func getTeamPosition(teamName string, scoreboard []database.Team) int {
	for index, item := range scoreboard {
		if item.Name == teamName {
			return index + 1
		}
	}

	return 0
}

func (h *Handler) SetDb(db *database.Database) {
	h.db = db
}

func (h *Handler) HandleGetScoreboard(w http.ResponseWriter, r *http.Request) {
	h.r.JSON(w, http.StatusOK, h.db.GetScoreboard())

}

func (h *Handler) HandleGetMatches(w http.ResponseWriter, r *http.Request) {
	h.r.JSON(w, http.StatusOK, h.db.GetMatches())
}

func (h *Handler) HandleGetTeams(w http.ResponseWriter, r *http.Request) {
	h.r.JSON(w, http.StatusOK, h.db.GetTeams())
}

func (h *Handler) HandleRenderTeamsPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"teams":       h.db.GetTeams(),
		"liveMatches": h.db.GetLiveMatches(),
		"scoreboard":  h.db.GetScoreboard(),
		"upcoming":    h.db.GetUpcomingMatches(),
		"week":        h.db.GetMatchWeekNumber(),
		"schedule":    h.db.GetPeacockSchedule(),
	}

	h.r.HTML(w, http.StatusOK, "index", data)
}
