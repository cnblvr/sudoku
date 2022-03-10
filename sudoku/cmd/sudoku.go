package main

import (
	sudoku "github.com/cnblvr/sudoku/sudoku/internal"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {

	// Initialize service sudoku
	srv, err := sudoku.NewService()
	if err != nil {
		log.Fatal().Err(err).Msg("sudoku.NewService failed")
	}

	r := mux.NewRouter()

	// Index page
	r.Path("/").Methods(http.MethodGet).HandlerFunc(srv.HandleIndex)

	// Websocket handle
	//r.Path("/ws").HandlerFunc()

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal().Err(err).Msg("http.ListenAndServe failed")
	}
}
