package user

import (
	"auth-server/internal/db"
	"auth-server/internal/security"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattn/go-sqlite3"
)

type userRepository interface {
	Create(email string, password string) error
	Login(email string, password string) (*User, error)
}

var userEmailTaken = errors.New("Email already taken")
var loginError = errors.New("Invalid email or password")

type defaultUserRepository struct {
	queries *db.Queries
}

func (r *defaultUserRepository) Create(email string, password string) error {
	hash, err := security.HashPassword(password)
	if err != nil {
		return fmt.Errorf("Error hashing password: %w", err)
	}

	_, err = r.queries.CreateUser(
		context.Background(),
		db.CreateUserParams{
			Email: email,
			Hash:  hash,
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
	dbUser, err := r.queries.FindUserByEmail(
		context.Background(),
		email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, loginError
		}
		return nil, err
	}

	ok, err := security.CheckPassword(password, dbUser.Hash)
	if err != nil {
		return nil, fmt.Errorf("Unable to login %v: Error comparing password: %w", email, err)
	}

	if !ok {
		return nil, loginError
	}

	user := User{
		ID:    strconv.Itoa(int(dbUser.ID)),
		Email: dbUser.Email,
	}

	return &user, nil
}
