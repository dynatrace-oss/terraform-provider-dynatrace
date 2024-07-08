package openpipeline

import (
	"context"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	openpipeline "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	caclib "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/openpipeline"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return service{}
}

type service struct {
	credentials *settings.Credentials
}

func (s *service) client() *caclib.Client {
	return nil
}

func (s service) List(ctx context.Context) (api.Stubs, error) {
	//TODO implement me
	return api.Stubs{}, nil
}

func (s service) Get(ctx context.Context, id string, v *openpipeline.Configuration) error {
	return nil
}

func (s service) SchemaID() string {
	return ""
}

func (s service) Create(ctx context.Context, v *openpipeline.Configuration) (*api.Stub, error) {
	return &api.Stub{}, nil
}

func (s service) Update(ctx context.Context, id string, v *openpipeline.Configuration) error {
	return nil
}

func (s service) Delete(ctx context.Context, id string) error {
	return nil
}
