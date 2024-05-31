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
	GetAllByUser(userID string) ([]App, error)
	Delete(clientID string, userID string) error
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

func (r *defaultAppRepository) GetAllByUser(userID string) ([]App, error) {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	apps, err := r.queries.GetAllAppsByUser(
		context.Background(),
		int64(userId),
	)
	if err != nil {
		return nil, err
	}

	var result []App
	for _, app := range apps {
		result = append(result, App{
			ID:           strconv.Itoa(int(app.ID)),
			ClientID:     app.Clientid.(string),
			ClientSecret: app.Hash,
			Name:         app.Name,
			RedirectURI:  app.Redirecturi,
			Type:         app.Type,
			UserID:       userID,
		})
	}

	return result, nil
}

func (r *defaultAppRepository) Delete(clientID string, userID string) error {
	userId, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}

	err = r.queries.DeleteApp(
		context.Background(),
		db.DeleteAppParams{
			Clientid: clientID,
			Userid:   int64(userId),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
