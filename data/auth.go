package data

// Auth is the authorization information contained in user's cookie and is used to control handlers and templates.
type Auth struct {
	// Flag indicates if the user is logged in.
	IsAuthorized bool `json:"isAuthorized"`
	// User ID.
	ID int64 `json:"id"`
	// TODO: token
}
