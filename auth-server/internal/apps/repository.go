package apps

import (
	"auth-server/internal/common"
	"auth-server/internal/db"
	"auth-server/internal/utils"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattn/go-sqlite3"
)

type appDBRegistry struct {
	common.Argon2idHash
	common.Timestamped

	ID          string
	ClientID    string
	Name        string
	RedirectURI string
	Type        string
	UserID      string
}

type registerInput struct {
	UserID      string
	Name        string
	Type        string
	RedirectURI string
}

type appRepository interface {
	Register(input registerInput) (*App, error)
}

var appEntries = []appDBRegistry{}

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
	if input.Type == "server-side" {
		clientSecret, err = utils.RandHexDecString(40)
		if err != nil {
			return nil, err
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
			Hash:        clientSecret,
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
