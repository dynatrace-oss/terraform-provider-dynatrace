package locations

import (
	"encoding/json"

	"github.com/dtcookie/dynatrace/rest"
	"github.com/dtcookie/dynatrace/rest/credentials"
)

// ServiceClient TODO: documentation
type ServiceClient struct {
	client *rest.Client
}

// NewService creates a new Service Client
// baseURL should look like this: "https://siz65484.live.dynatrace.com/api/config/v1"
// token is an API Token
func NewService(baseURL string, token string) *ServiceClient {
	credentials := credentials.New(token)
	config := rest.Config{}
	client := rest.NewClient(&config, baseURL, credentials)

	return &ServiceClient{client: client}
}

// // Get TODO: documentation
// func (cs *ServiceClient) Get(id string) (*Credentials, error) {
// 	if len(id) == 0 {
// 		return nil, errors.New("empty ID provided for the configuration to fetch")
// 	}

// 	var err error
// 	var bytes []byte

// 	if bytes, err = cs.client.GET(fmt.Sprintf("/credentials/%s", id), 200); err != nil {
// 		return nil, err
// 	}
// 	var credentials Credentials
// 	if err = json.Unmarshal(bytes, &credentials); err != nil {
// 		return nil, err
// 	}
// 	return &credentials, nil
// }

// ListAll TODO: documentation
func (cs *ServiceClient) List() (*SyntheticLocations, error) {
	var err error
	var bytes []byte

	rest.Verbose = true

	if bytes, err = cs.client.GET("/synthetic/locations", 200); err != nil {
		return nil, err
	}
	var stubList SyntheticLocations
	if err = json.Unmarshal(bytes, &stubList); err != nil {
		return nil, err
	}
	return &stubList, nil
}
