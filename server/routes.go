package server

import "net/http"

func (s *Server) Routes() {
	s.staticFiles()

	s.router.Get("/api/scoreboard", s.handler.HandleGetScoreboard)
	s.router.Get("/api/matches", s.handler.HandleGetMatches)
	s.router.Get("/api/teams", s.handler.HandleGetTeams)

	s.router.Get("/teams", s.handler.HandleRenderTeamsPage)
}

func (s *Server) staticFiles() {
	fs := http.FileServer(http.Dir("./templates/public"))
	s.router.Handle("/public/*", http.StripPrefix("/public/", fs))
}
