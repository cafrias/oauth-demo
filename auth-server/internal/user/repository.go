package user

import (
	"auth-server/internal/common"
	"auth-server/internal/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattn/go-sqlite3"
)

type userDBRegistry struct {
	common.Argon2idHash
	common.Timestamped

	ID    string
	Email string
}

type userRepository interface {
	Create(email string, password string) error
	Login(email string, password string) (*User, error)
}

var entries = []userDBRegistry{
	{
		ID:    "1",
		Email: "a@a.com",
		Argon2idHash: common.Argon2idHash{
			Hash: "1234",
		},
	},
}

var userEmailTaken = errors.New("Email already taken")
var loginError = errors.New("Invalid email or password")

type defaultUserRepository struct {
	queries *db.Queries
}

func (r *defaultUserRepository) Create(email string, password string) error {
	_, err := r.queries.CreateUser(
		context.Background(),
		db.CreateUserParams{
			Email: email,
			Hash:  password,
		},
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return fmt.Errorf("Error creating user '%s': %w", email, userEmailTaken)
			}
		}
		return fmt.Errorf("Error creating user '%s': %w", email, err)
	}

	return nil
}

func (r *defaultUserRepository) Login(email string, password string) (*User, error) {
	dbUser, err := r.queries.Login(
		context.Background(),
		db.LoginParams{
			Email: email,
			// TODO: password must be hashed
			Hash: password,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, loginError
		}
		return nil, err
	}

	user := User{
		ID:    strconv.Itoa(int(dbUser.ID)),
		Email: dbUser.Email,
	}

	return &user, nil
}
