package data

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"io"
)

const (
	// EndpointIndex is a path to the main page of the site.
	EndpointIndex = "/"
	// EndpointLogin is a path to the login page.
	EndpointLogin = "/login"
	// EndpointSignup is a path to the registration page.
	EndpointSignup = "/signup"
	// EndpointLogout is a path to the logout handler.
	EndpointLogout = "/logout"
	// EndpointUserInfo is a path to the user's info page.
	EndpointUserInfo = "/info"
	// EndpointSudokuPlay is a path to the puzzle generator page/handler.
	EndpointSudokuPlay = "/sudoku/play"
	endpointSudokuGame = "/sudoku/%s"
)

func EndpointSudoku(sudokuID string) string {
	return fmt.Sprintf(endpointSudokuGame, sudokuID)
}

func GenerateReqID() string {
	buf := make([]byte, 6)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "0000"
	}
	return fmt.Sprintf("%x-%x", buf[0:3], buf[3:6])
}

func GetCtxReqID(ctx context.Context) string {
	val := ctx.Value("reqid")
	if val == nil {
		return ""
	}
	reqid, ok := val.(string)
	if !ok {
		return ""
	}
	return reqid
}

func SetCtxReqID(ctx context.Context, reqid string) context.Context {
	return context.WithValue(ctx, "reqid", reqid)
}

func GetCtxLog(ctx context.Context) zerolog.Logger {
	val := ctx.Value("logger")
	if val == nil {
		return zlog.Logger
	}
	logger, ok := val.(zerolog.Logger)
	if !ok {
		return zlog.Logger
	}
	return logger
}

func SetCtxLog(ctx context.Context, logger zerolog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}
