package user

import (
	"errors"
	"fmt"
)

type user struct {
	ID    string
	Email string
}

type userDBRegistry struct {
	ID      string
	Email   string
	Hash    string
	Salt    string
	Created string
	Updated string
	Deleted bool
}

type userRepository interface {
	FindByEmail(email string) (*user, error)
	Create(email string, password string) error
	Save(user user) error
}

var entries = []userDBRegistry{
	{
		ID:    "1",
		Hash:  "1234",
		Email: "a@a.com",
	},
}

var userNotFound = errors.New("User not found")
var userEmailTaken = errors.New("Email already taken")

type defaultUserRepository struct{}

func (r *defaultUserRepository) FindByEmail(email string) (*user, error) {
	for _, entry := range entries {
		if entry.Email == email {
			return &user{
				ID:    entry.ID,
				Email: entry.Email,
			}, nil
		}
	}
	return nil, fmt.Errorf("Error finding user '%s': %w", email, userNotFound)
}

func (r *defaultUserRepository) Save(user user) error {
	entries = append(entries, userDBRegistry{
		ID:    user.ID,
		Email: user.Email,
	})
	return nil
}

func (r *defaultUserRepository) Create(email string, password string) error {
	for _, entry := range entries {
		if entry.Email == email {
			return fmt.Errorf("Error creating user '%s': %w", email, userEmailTaken)
		}
	}

	entries = append(entries, userDBRegistry{
		Email: email,
		Hash:  password,
	})

	return nil
}
