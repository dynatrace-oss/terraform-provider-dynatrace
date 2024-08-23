package openpipeline

import (
	"context"
	"encoding/json"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	openpipeline "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/openpipeline/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace/dynatrace-configuration-as-code-core/clients"
	caclib "github.com/dynatrace/dynatrace-configuration-as-code-core/clients/openpipeline"
	"golang.org/x/oauth2/clientcredentials"
)

func Service(credentials *settings.Credentials) settings.CRUDService[*openpipeline.Configuration] {
	return &service{credentials}
}

type service struct {
	credentials *settings.Credentials
}

func (s *service) client() *caclib.Client {
	factory := clients.Factory().
		WithUserAgent("Dynatrace Terraform Provider").
		WithPlatformURL(s.credentials.Automation.EnvironmentURL).
		WithOAuthCredentials(clientcredentials.Config{
			ClientID:     s.credentials.Automation.ClientID,
			ClientSecret: s.credentials.Automation.ClientSecret,
			TokenURL:     s.credentials.Automation.TokenURL,
		})

	openPipelineClient, _ := factory.OpenPipelineClient()
	return openPipelineClient

}

func (s service) List(ctx context.Context) (api.Stubs, error) {
	//TODO implement me
	return api.Stubs{}, nil
}

func (s service) Get(ctx context.Context, id string, v *openpipeline.Configuration) error {
	response, err := s.client().Get(ctx, id)
	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return rest.Envelope(response.Data, s.credentials.Automation.EnvironmentURL, "GET")
	}

	return json.Unmarshal(response.Data, &v)
}

func (s service) SchemaID() string {
	return "platform:openpipeline"
}

func (s service) Create(_ context.Context, v *openpipeline.Configuration) (*api.Stub, error) {
	return &api.Stub{ID: v.Kind}, nil
}

func (s service) Update(_ context.Context, id string, v *openpipeline.Configuration) error {
	return nil
}

func (s service) Delete(_ context.Context, id string) error {
	return nil
}
