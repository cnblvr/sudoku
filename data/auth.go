package data

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

// Auth is the authorization information contained in user's cookie and is used to control handlers and templates.
type Auth struct {
	// Flag indicates if the user is logged in.
	IsAuthorized bool `json:"-"`
	// User ID.
	ID    int64     `json:"-"`
	Token uuid.UUID `json:"token"`
}

func (a Auth) UserID() int64 {
	return a.ID
}

func (a Auth) TokenID() uuid.UUID {
	return a.Token
}

func (a Auth) IsEmptyTokenID() bool {
	return a.Token.String() == "00000000-0000-0000-0000-000000000000"
}

func GetCtxAuth(ctx context.Context) *Auth {
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

func SetCtxAuth(ctx context.Context, auth *Auth) context.Context {
	return context.WithValue(ctx, "auth", &auth)
}
