package repository_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/repository"
	"testing"
	"text/template"
	"time"
)

func TestUserRepository(t *testing.T) {
	userRepository := data.UserRepository(repository.New(
		flushDB(t, newRedisTestPool(t))),
	)
	ctx := context.Background()

	const (
		wantUserID1       = 1
		wantUserUsername1 = "username1"
		wantUserID2       = 2
		wantUserUsername2 = "username2"
	)
	var (
		wantUserHash1 = []byte{'h', 'a', 's', 'h', '1'}
		wantUserSalt1 = []byte{'s', 'a', 'l', 't', '1'}
	)

	user1, err := userRepository.CreateUser(ctx, wantUserUsername1)
	if err != nil {
		t.Fatal(err)
	}
	if user1.ID() != wantUserID1 {
		t.Fatalf("user1's id is not '%d'", wantUserID1)
	}
	if user1.Username != wantUserUsername1 {
		t.Fatalf("user1's username is not '%s'", wantUserUsername1)
	}

	user2, err := userRepository.CreateUser(ctx, wantUserUsername2)
	if err != nil {
		t.Fatal(err)
	}
	if user2.ID() != wantUserID2 {
		t.Fatalf("user2's id is not '%d'", wantUserID2)
	}
	if user2.Username != wantUserUsername2 {
		t.Fatalf("user2's username is not '%s'", wantUserUsername2)
	}

	user1.PasswordHash = wantUserHash1
	user1.PasswordSalt = wantUserSalt1
	if err := userRepository.UpdateUser(ctx, user1); err != nil {
		t.Fatal(err)
	}

	user1Updated, err := userRepository.GetUserByID(ctx, user1.ID())
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(user1Updated.PasswordHash, wantUserHash1) {
		t.Fatalf("user1's hash not '%s'", string(wantUserHash1))
	}
	if !bytes.Equal(user1Updated.PasswordSalt, wantUserSalt1) {
		t.Fatalf("user1's salt not '%s'", string(wantUserSalt1))
	}

	user2.Username = wantUserUsername1
	if err := userRepository.UpdateUser(ctx, user2); err == nil {
		t.Fatalf("user2's error is nil")
	}

	const wantUserUsername2Changed = wantUserUsername2 + "_changed"
	user2.Username = wantUserUsername2Changed
	if err := userRepository.UpdateUser(ctx, user2); err != nil {
		t.Fatal(err)
	}
	user2Updated, err := userRepository.GetUserByUsername(ctx, wantUserUsername2Changed)
	if err != nil {
		t.Fatal(err)
	}
	if user2Updated.Username != wantUserUsername2Changed {
		t.Fatalf("user2Changed's username not '%s'", string(wantUserUsername2Changed))
	}

}

func TestReadOnlyFieldsInTemplates(t *testing.T) {
	userRepository := data.UserRepository(repository.New(
		flushDB(t, newRedisTestPool(t))),
	)
	ctx := context.Background()

	user, err := userRepository.CreateUser(ctx, "username")
	if err != nil {
		t.Fatal(err)
	}

	gotBuf := bytes.NewBuffer(nil)
	if err := template.Must(template.New("").Parse(
		`User: '{{.Username}}'; ID: {{.ID}}; CreatedAt: '{{.CreatedAt.Format "2006-01-02T15"}}'; UpdatedAt: {{.UpdatedAt.Day}}`,
	)).Execute(gotBuf, user); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	var want = fmt.Sprintf(`User: 'username'; ID: 1; CreatedAt: '%s'; UpdatedAt: %d`, now.Format("2006-01-02T15"), now.Day())
	if got := gotBuf.String(); got != want {
		t.Fatalf("want `%s`, got `%s`", got, want)
	}
}
