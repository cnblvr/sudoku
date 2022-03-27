package main

import (
	generator "github.com/cnblvr/sudoku/generator/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Logger initialization
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log.Logger = zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "2006-01-02 15:04:05.000000Z"
	})).With().Timestamp().Caller().Logger()

	// Initialize Generator service
	srv, err := generator.NewService()
	if err != nil {
		log.Fatal().Err(err).Msg("sudoku.NewService failed")
	}

	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("generator.Service failed")
	}
}
