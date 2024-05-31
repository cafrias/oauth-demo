package apps

import (
	"auth-server/internal/db"
	"auth-server/internal/security"
	"auth-server/internal/utils"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattn/go-sqlite3"
)

type registerInput struct {
	UserID      string
	Name        string
	Type        string
	RedirectURI string
}

type appRepository interface {
	Register(input registerInput) (*App, error)
}

type defaultAppRepository struct {
	queries *db.Queries
}

var clientIDTaken = errors.New("Client ID already taken")

func (r *defaultAppRepository) Register(input registerInput) (*App, error) {
	userId, err := strconv.Atoi(input.UserID)
	if err != nil {
		return nil, err
	}

	clientID, err := utils.RandHexDecString(20)
	if err != nil {
		return nil, err
	}

	// TODO: check that clientID is unique

	var clientSecret string
	var hash string
	if input.Type == "server-side" {
		clientSecret, err = utils.RandHexDecString(40)
		if err != nil {
			return nil, err
		}

		hash, err = security.HashPassword(clientSecret)
		if err != nil {
			return nil, fmt.Errorf("Error hashing client secret: %w", err)
		}
	}

	id, err := r.queries.CreateApp(
		context.Background(),
		db.CreateAppParams{
			Name:        input.Name,
			Userid:      int64(userId),
			Type:        input.Type,
			Redirecturi: input.RedirectURI,
			Clientid:    clientID,
			Hash:        hash,
		},
	)
	if err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return nil, fmt.Errorf("Error creating app: %w", clientIDTaken)
		}

		return nil, err
	}

	return &App{
		ID:           strconv.Itoa(int(id)),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Name:         input.Name,
		RedirectURI:  input.RedirectURI,
		Type:         input.Type,
		UserID:       input.UserID,
	}, nil

}
