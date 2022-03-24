package data

import (
	"context"
	"fmt"
)

var (
	ErrUsernameIsBusy = fmt.Errorf("username is busy")
	ErrUserNotFound   = fmt.Errorf("user not found")
)

type UserRepository interface {
	// Errors: ErrUsernameIsBusy, unknown.
	CreateUser(ctx context.Context, username string) (*User, error)
	// Errors: ErrUserNotFound, unknown.
	GetUserByID(ctx context.Context, id int64) (*User, error)
	// Errors: ErrUserNotFound, unknown.
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// Errors: ErrUserNotFound, unknown.
	UpdateUser(ctx context.Context, user *User) error
	// Errors: ErrUserNotFound, unknown.
	DeleteUser(ctx context.Context, id int64) error
}

type User struct {
	IDint64Getter   `json:"id"`
	Username        string `json:"username"`
	PasswordHash    []byte `json:"password_hash,omitempty"`
	PasswordSalt    []byte `json:"password_salt,omitempty"`
	CreatedAtGetter `json:"created_at"`
	UpdatedAtGetter `json:"updated_at"`
}
