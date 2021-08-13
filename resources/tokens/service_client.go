package tokens

import (
	"encoding/json"
	"fmt"

	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/rest/credentials"
)

// Service TODO: documentation
type Service struct {
	client *rest.Client
}

// NewService TODO: documentation
// "https://#######.live.dynatrace.com/api/config/v1", "###########"
func NewServiceRest(baseURL string, token string) *Service {
	rest.Verbose = false
	credentials := credentials.New(token)
	config := rest.Config{}
	client := rest.NewClient(&config, baseURL, credentials)

	return &Service{client: client}
}

// Create TODO: documentation
func (cs *Service) Create(token *APIToken) (*APIToken, error) {
	var err error
	var bytes []byte

	if bytes, err = cs.client.POST("/apiTokens", token, 201); err != nil {
		return nil, err
	}
	var resultToken APIToken
	if err = json.Unmarshal(bytes, &resultToken); err != nil {
		return nil, err
	}
	resultToken.Scopes = token.Scopes
	return &resultToken, nil
}

// Update TODO: documentation
func (cs *Service) Update(token *APIToken) error {
	// if _, err := cs.client.PUT(fmt.Sprintf("/alertingProfiles/%s", *token.ID), token, 204); err != nil {
	// 	return err
	// }
	return nil
}

// Delete TODO: documentation
func (cs *Service) Delete(id string) error {
	if _, err := cs.client.DELETE(fmt.Sprintf("/apiTokens/%s", id), 204); err != nil {
		return err
	}
	return nil
}

// Get TODO: documentation
func (cs *Service) Get(id string) (*APIToken, error) {
	var err error
	var bytes []byte

	if bytes, err = cs.client.GET(fmt.Sprintf("/apiTokens/%s", id), 200); err != nil {
		return nil, err
	}
	var token APIToken
	if err = json.Unmarshal(bytes, &token); err != nil {
		return nil, err
	}
	return &token, nil
}
