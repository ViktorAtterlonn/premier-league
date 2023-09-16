package handlers

import (
	"net/http"
	"scraper/database"

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
		}),
	}
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
	h.r.HTML(w, http.StatusOK, "index", h.db.GetTeams())
}
