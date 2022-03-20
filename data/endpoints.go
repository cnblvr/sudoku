package data

import "fmt"

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
