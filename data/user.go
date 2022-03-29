package data

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	ErrUsernameIsBusy = fmt.Errorf("username is busy")
	ErrUserNotFound   = fmt.Errorf("user not found")
	ErrTokenNotFound  = fmt.Errorf("token not found")
)

type UserRepository interface {
	CreateToken(ctx context.Context, userID int64, expiration time.Duration) (*TokenInfo, error)
	GetToken(ctx context.Context, id uuid.UUID) (*TokenInfo, error)
	DeleteToken(ctx context.Context, id uuid.UUID) error
	// Errors: ErrUsernameIsBusy, unknown.
	CreateUser(ctx context.Context, username string) (*User, error)
	// Errors: ErrUserNotFound, unknown.
	GetUserByID(ctx context.Context, id int64) (*User, error)
	// Errors: ErrUserNotFound, unknown.
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// Errors: ErrUserNotFound, ErrUsernameIsBusy, unknown.
	UpdateUser(ctx context.Context, user *User) error
	// Errors: ErrUserNotFound, unknown.
	DeleteUser(ctx context.Context, id int64) error
}

type User struct {
	ID           int64    `json:"id"`
	Username     string   `json:"username"`
	PasswordHash []byte   `json:"password_hash,omitempty"`
	PasswordSalt []byte   `json:"password_salt,omitempty"`
	CreatedAt    DateTime `json:"created_at"`
	UpdatedAt    DateTime `json:"updated_at"`
}

type TokenInfo struct {
	ID         uuid.UUID     `json:"id"`
	Expiration time.Duration `json:"expiration"`
	UserID     int64         `json:"user_id,omitempty"`
}

type UserOrToken interface {
	UserID() int64
	TokenID() uuid.UUID
	IsEmptyTokenID() bool
}

const DefaultTokenExpiration = time.Hour

func GetCtxUser(ctx context.Context) *User {
	val := ctx.Value("user")
	if val == nil {
		return &User{}
	}
	user, ok := val.(*User)
	if !ok {
		return &User{}
	}
	return user
}

func SetCtxUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, "user", &user)
}
