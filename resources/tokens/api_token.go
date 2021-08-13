package tokens

import (
	"github.com/dtcookie/hcl"
)

type APIToken struct {
	ID                  *string  `json:"id"`
	Name                string   `json:"name"`
	Enabled             bool     `json:"enabled"`
	PersonalAccessToken bool     `json:"personalAccessToken"`
	ExpirationDate      string   `json:"expirationDate,omitempty"`
	Scopes              []string `json:"scopes"`
	Token               *string  `json:"token,omitempty"`
}

func (me *APIToken) Schema() map[string]*hcl.Schema {
	return map[string]*hcl.Schema{
		"name": {
			Type:        hcl.TypeString,
			Description: "The name of the token",
			Required:    true,
		},
		"enabled": {
			Type:        hcl.TypeBool,
			Description: "This API Token is enabled or disabled",
			Optional:    true,
		},
		"scopes": {
			Type:        hcl.TypeSet,
			Description: "The scopes of the API Token",
			Required:    true,
			Elem:        &hcl.Schema{Type: hcl.TypeString},
		},
		"expiration_date": {
			Type:        hcl.TypeString,
			Description: "Expiration Date",
			Optional:    true,
		},
		"token": {
			Type:        hcl.TypeString,
			Description: "The token is being generated upon creation of this resource. You cannot specify it nor can you change it.",
			Computed:    true,
			Sensitive:   true,
		},
	}
}

func (me *APIToken) MarshalHCL() (map[string]interface{}, error) {
	properties := hcl.Properties{}
	if _, err := properties.EncodeAll(map[string]interface{}{
		"name":            me.Name,
		"expiration_date": me.ExpirationDate,
		"scopes":          me.Scopes,
	}); err != nil {
		return nil, err
	}
	if me.Token != nil {
		if err := properties.Encode("token", me.Token); err != nil {
			return nil, err
		}
	}
	return properties, nil
}

func (me *APIToken) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("scopes", &me.Scopes); err != nil {
		return err
	}
	if err := decoder.Decode("name", &me.Name); err != nil {
		return err
	}
	// if err := decoder.Decode("token", &me.Token); err != nil {
	// 	return err
	// }
	if err := decoder.Decode("expiration_date", &me.ExpirationDate); err != nil {
		return err
	}
	return nil
}
