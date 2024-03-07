package monitoring_config

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	monitoring_config "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/extensions/monitoring_config/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*monitoring_config.Settings] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (me *service) Get(id string, v *monitoring_config.Settings) error {
	// client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)
	return nil
}

func (me *service) SchemaID() string {
	return "v2:extensions:monitoring:config"
}

func (me *service) List() (api.Stubs, error) {
	var stubs api.Stubs
	// client := rest.DefaultClient(me.credentials.URL, me.credentials.Token)

	return stubs, nil
}

func (me *service) Validate(v *monitoring_config.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(v *monitoring_config.Settings) (*api.Stub, error) {
	var id string
	var name string
	return &api.Stub{ID: id, Name: name}, nil
}

func (me *service) Update(id string, v *monitoring_config.Settings) error {
	return nil
}

func (me *service) Delete(id string) error {
	return nil // no endpoint for that
}

func (me *service) New() *monitoring_config.Settings {
	return new(monitoring_config.Settings)
}

func (me *service) Name() string {
	return me.SchemaID()
}
