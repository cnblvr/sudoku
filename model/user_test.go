package model

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestUser(t *testing.T) {
	conn, err := redis.Dial("tcp", "localhost:6379", redis.DialDatabase(1))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Do("FLUSHDB", "SYNC"); err != nil {
		t.Fatal(err)
	}

	const wantID, wantUsername = int64(1), "username"
	var wantHash, wantSalt = []byte{0x00, 0x01, 0x02, 0xff}, []byte{0x82, 0x88, 0x00, 0xf3, 0x19}
	if _, err := conn.Do("SET", fmt.Sprintf("user:%d", wantID), wantUsername); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Do("SET", keyLastUserID(), 1); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Do("SET", keyUsername(wantUsername), wantID); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Do("SET", keyUserPasswordHash(wantID), wantHash); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Do("SET", keyUserPasswordSalt(wantID), wantSalt); err != nil {
		t.Fatal(err)
	}

	var user User
	t.Run("UserByID()", func(t *testing.T) {
		u, isExists, err := UserByID(conn, wantID)
		if err != nil {
			t.Fatal(err)
		}
		if !isExists {
			t.Fatalf("user %d not found", wantID)
		}
		user = u
	})

	t.Run("UserByUsername()", func(t *testing.T) {
		userByUsername, isExists, err := UserByUsername(conn, wantUsername)
		if err != nil {
			t.Fatal(err)
		}
		if !isExists {
			t.Fatalf("user %s not found", wantUsername)
		}
		if username, err := userByUsername.Username(); err != nil {
			t.Fatal(err)
		} else if username != wantUsername {
			t.Fatalf("want username is not %s. Got: %s", wantUsername, username)
		}
	})

	t.Run("User.ID()", func(t *testing.T) {
		if id := user.ID(); id != wantID {
			t.Fatalf("want ID is not %d. Got %d", wantID, id)
		}
	})

	t.Run("User.Username()", func(t *testing.T) {
		if username, err := user.Username(); err != nil {
			t.Fatal(err)
		} else if username != wantUsername {
			t.Fatalf("want username is not %s. Got: %s", wantUsername, username)
		}
	})

	t.Run("User.SetUsername()", func(t *testing.T) {
		const wantUsernameNew = wantUsername + "new"
		if err := user.SetUsername(wantUsernameNew); err != nil {
			t.Fatal(err)
		}
		if usernameNew, err := user.Username(); err != nil {
			t.Fatal(err)
		} else if usernameNew != wantUsernameNew {
			t.Fatalf("want username is not %s. Got: %s", wantUsernameNew, usernameNew)
		}
		if err := user.SetUsername(wantUsername); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("User.PasswordHash()", func(t *testing.T) {
		hash, err := user.PasswordHash()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(hash, wantHash) {
			t.Fatalf("want hash is not %s. Got: %s",
				base64.StdEncoding.EncodeToString(hash),
				base64.StdEncoding.EncodeToString(wantHash),
			)
		}
	})

	t.Run("User.SetPasswordHash()", func(t *testing.T) {
		var wantHashNew = append(wantHash, 0x83, 0xb2, 0x11)
		if err := user.SetPasswordHash(wantHashNew); err != nil {
			t.Fatal(err)
		}
		hash, err := user.PasswordHash()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(hash, wantHashNew) {
			t.Fatalf("want hash is not %s. Got: %s",
				base64.StdEncoding.EncodeToString(hash),
				base64.StdEncoding.EncodeToString(wantHashNew),
			)
		}
		if err := user.SetPasswordHash(wantHash); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("User.PasswordSalt()", func(t *testing.T) {
		salt, err := user.PasswordSalt()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(salt, wantSalt) {
			t.Fatalf("want salt is not %s. Got: %s",
				base64.StdEncoding.EncodeToString(salt),
				base64.StdEncoding.EncodeToString(wantSalt),
			)
		}
	})

	t.Run("User.SetPasswordSalt()", func(t *testing.T) {
		var wantSaltNew = append(wantHash, 0x93, 0x5f, 0x00, 0xdd)
		if err := user.SetPasswordSalt(wantSaltNew); err != nil {
			t.Fatal(err)
		}
		salt, err := user.PasswordSalt()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(salt, wantSaltNew) {
			t.Fatalf("want salt is not %s. Got: %s",
				base64.StdEncoding.EncodeToString(salt),
				base64.StdEncoding.EncodeToString(wantSaltNew),
			)
		}
		if err := user.SetPasswordSalt(wantSalt); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("IsUsernameVacant()", func(t *testing.T) {
		isVacant, err := IsUsernameVacant(conn, wantUsername)
		if err != nil {
			t.Fatal(err)
		}
		if isVacant {
			t.Fatalf("username %s is vacant", wantUsername)
		}
		const wantUsernameVacant = wantUsername + "Vacant"
		isVacant, err = IsUsernameVacant(conn, wantUsernameVacant)
		if err != nil {
			t.Fatal(err)
		}
		if !isVacant {
			t.Fatalf("username %s is not vacant", wantUsernameVacant)
		}
	})

	t.Run("NewUser()", func(t *testing.T) {
		const wantUsernameNew = "username2"
		wantIDNew, err := redis.Int64(conn.Do("GET", keyLastUserID()))
		if err != nil {
			t.Fatal(err)
		}
		wantIDNew += 1
		userNew, err := NewUser(conn, wantUsernameNew)
		if err != nil {
			t.Fatal(err)
		}
		if username, err := userNew.Username(); err != nil {
			t.Fatal(err)
		} else if username != wantUsernameNew {
			t.Fatalf("want username is not %s. Got: %s", wantUsernameNew, username)
		}
		if id := userNew.ID(); id != wantIDNew {
			t.Fatalf("want ID is not %d. Got: %d", wantIDNew, id)
		}

		createdAt, err := userNew.CreatedAt()
		if err != nil {
			t.Fatal(err)
		}
		if createdAt.IsZero() {
			t.Fatalf("CreatedAt is zero")
		}
	})
}
