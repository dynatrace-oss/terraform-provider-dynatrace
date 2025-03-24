package naming_order

import (
	"context"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/order/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *rest.Credentials, schemaID string, staticID string, staticName string, basePath string) settings.CRUDService[*order.Settings] {
	return &service{client: rest.APITokenClient(credentials), schemaID: schemaID, StaticID: staticID, StaticName: staticName, BasePath: basePath}
}

type service struct {
	client     rest.Client
	schemaID   string
	StaticID   string
	StaticName string
	BasePath   string
}

func (s *service) Create(ctx context.Context, v *order.Settings) (*api.Stub, error) {
	// Creating this resource doesn't really add a setting on the remote side
	// It just defines a specific order
	return &api.Stub{ID: s.StaticID, Name: s.StaticName}, s.Update(ctx, "", v)
}

func (s *service) Get(ctx context.Context, id string, v *order.Settings) (err error) {
	cfg := ctx.Value(settings.ContextKeyStateConfig)
	var currentIDs []string
	if orderConfig, ok := cfg.(*order.Settings); ok && len(orderConfig.IDs) > 0 {
		currentIDs = orderConfig.IDs
	}

	var listResponse struct {
		Values []struct {
			ID string `json:"id"`
		} `json:"values"`
	}
	req := s.client.Get(ctx, s.BasePath, 200)
	if err = req.Finish(&listResponse); err != nil {
		return err
	}
	for _, entry := range listResponse.Values {
		if len(currentIDs) == 0 || slices.Contains(currentIDs, entry.ID) {
			v.IDs = append(v.IDs, entry.ID)
		}
	}
	return nil
}

type ListResponse struct {
	Values []struct {
		ID string `json:"id"`
	} `json:"values"`
}

func (s *service) List(ctx context.Context) (stubs api.Stubs, err error) {
	return append(stubs, &api.Stub{ID: s.StaticID, Name: s.StaticName}), nil
}

func (s *service) SchemaID() string {
	return s.schemaID
}

func (s *service) Update(ctx context.Context, id string, v *order.Settings) (err error) {
	var payload = struct {
		Values []struct {
			ID string `json:"id"`
		} `json:"values"`
	}{}
	for _, id := range v.IDs {
		payload.Values = append(payload.Values, struct {
			ID string `json:"id"`
		}{ID: id})
	}
	return s.client.Put(ctx, s.BasePath+"/order", &payload, 204).Finish()
}

func (s *service) Delete(ctx context.Context, id string) error {
	// Deleting this resource does not actually destroy anything on the remote side
	// We simply forget about the order of things
	return nil
}
