package data

import "context"

// Auth is the authorization information contained in user's cookie and is used to control handlers and templates.
type Auth struct {
	// Flag indicates if the user is logged in.
	IsAuthorized bool `json:"isAuthorized"`
	// User ID.
	ID int64 `json:"id"`
	// TODO: token
}

func GetAuth(ctx context.Context) *Auth {
	val := ctx.Value("auth")
	if val == nil {
		return &Auth{}
	}
	auth, ok := val.(*Auth)
	if !ok {
		return &Auth{}
	}
	return auth
}
