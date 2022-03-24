package main

import (
	sudoku "github.com/cnblvr/sudoku/sudoku/internal"
	"github.com/cnblvr/sudoku/sudoku/static"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// Logger initialization
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log.Logger = zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "2006-01-02 15:04:05.000000Z"
	})).With().Timestamp().Caller().Logger()

	// Initialize Sudoku service
	srv, err := sudoku.NewService()
	if err != nil {
		log.Fatal().Err(err).Msg("sudoku.NewService failed")
	}

	r := mux.NewRouter()

	// Router for static files: JS, CSS, images.
	rStatic := r.NewRoute().Subrouter()
	rStatic.Path("/favicon.ico").Methods(http.MethodGet).Handler(http.FileServer(http.FS(static.Favicon)))
	rStatic.PathPrefix("/css").Methods(http.MethodGet).Handler(http.FileServer(http.FS(static.CSS)))
	rStatic.PathPrefix("/js").Methods(http.MethodGet).Handler(http.FileServer(http.FS(static.JS)))

	// Router rPages for pages or actions that do not require authorization: index, login/signup, anonymous game.
	rPages := r.NewRoute().Subrouter()
	rPages.Use(srv.MiddlewareCookies)
	// Router rAuth manages pages or actions for an authorized users: info about last games, logout.
	rAuth := r.NewRoute().Subrouter()
	rAuth.Use(srv.MiddlewareCookies, srv.MiddlewareMustBeAuthorized)

	// Main page
	rPages.Path("/").Methods(http.MethodGet).HandlerFunc(srv.HandleIndex)
	// Login page and handler
	rPages.Path("/login").Methods(http.MethodGet, http.MethodPost).HandlerFunc(srv.HandleLogin)
	// Logout handler
	rPages.Path("/logout").Methods(http.MethodGet).HandlerFunc(srv.HandleLogout)
	// Signup page and handler
	rPages.Path("/signup").Methods(http.MethodGet, http.MethodPost).HandlerFunc(srv.HandleSignup)
	// User's info page and handler
	//rAuth.Path("/info").Methods(http.MethodGet, http.MethodPost).HandlerFunc(srv.HandleUserInfo)
	// Puzzle create game page
	rPages.Path("/sudoku/play").Methods(http.MethodGet).HandlerFunc(srv.HandleSudokuCreate)
	// Puzzle page
	rPages.Path("/sudoku/{game_id}").Methods(http.MethodGet).HandlerFunc(srv.HandleSudoku)

	// Websocket handler
	rPages.Path("/ws").Methods(http.MethodGet).HandlerFunc(srv.HandleWebsocket)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal().Err(err).Msg("http.ListenAndServe failed")
	}
}
