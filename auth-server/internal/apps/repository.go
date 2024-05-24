package apps

import "auth-server/internal/utils"

type appDBRegistry struct {
	ID          string
	ClientID    string
	Hash        string
	Salt        string
	Name        string
	RedirectURI string
	Type        string
	UserID      string
	Created     string
	Updated     string
	Deleted     bool
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

type defaultAppRepository struct{}

func (r *defaultAppRepository) Register(input registerInput) (*App, error) {
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

	reg := appDBRegistry{
		ClientID:    clientID,
		Hash:        clientSecret,
		Name:        input.Name,
		RedirectURI: input.RedirectURI,
		Type:        input.Type,
		UserID:      input.UserID,
	}

	appEntries = append(appEntries, reg)

	return &App{
		ID:           reg.ID,
		ClientID:     reg.ClientID,
		ClientSecret: clientSecret,
		Name:         reg.Name,
		RedirectURI:  reg.RedirectURI,
		Type:         reg.Type,
		UserID:       reg.UserID,
	}, nil

}
