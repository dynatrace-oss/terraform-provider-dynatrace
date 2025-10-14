package boundaries

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	boundaries "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/boundaries/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/clean"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*boundaries.PolicyBoundary] {
	return &BoundaryServiceClient{
		clientID:     credentials.IAM.ClientID,
		accountID:    credentials.IAM.AccountID,
		clientSecret: credentials.IAM.ClientSecret,
		tokenURL:     credentials.IAM.TokenURL,
		endpointURL:  credentials.IAM.EndpointURL,
	}
}

type BoundaryServiceClient struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *BoundaryServiceClient) ClientID() string {
	return me.clientID
}

func (me *BoundaryServiceClient) AccountID() string {
	return me.accountID
}

func (me *BoundaryServiceClient) ClientSecret() string {
	return me.clientSecret
}

func (me *BoundaryServiceClient) TokenURL() string {
	return me.tokenURL
}

func (me *BoundaryServiceClient) EndpointURL() string {
	return me.endpointURL
}

func (me *BoundaryServiceClient) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(ctx, fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")), 200, false); err != nil {
		return nil, err
	}

	var response ListPolicyBoundariesResponse
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	if len(response.PolicyBoundaries) == 0 {
		return api.Stubs{}, nil
	}
	stubs := api.Stubs{}
	for _, boundary := range response.PolicyBoundaries {
		stubs = append(stubs, &api.Stub{ID: boundary.UUID, Name: boundary.Name})
	}
	return stubs, nil
}

func (me *BoundaryServiceClient) Get(ctx context.Context, id string, v *boundaries.PolicyBoundary) error {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.GET(ctx, fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id), 200, false); err != nil {
		return err
	}

	if err = json.Unmarshal(responseBytes, v); err != nil {
		return err
	}
	return nil
}

func (me *BoundaryServiceClient) SchemaID() string {
	return "accounts:iam:boundaries"
}

func (me *BoundaryServiceClient) Create(ctx context.Context, v *boundaries.PolicyBoundary) (*api.Stub, error) {
	var err error
	var responseBytes []byte

	client := iam.NewIAMClient(me)

	if responseBytes, err = client.POST(
		ctx,
		fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:")),
		v,
		201,
		false,
	); err != nil {
		return nil, err
	}

	response := struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	}{}

	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}

	return &api.Stub{ID: response.UUID, Name: response.Name}, nil
}

func (me *BoundaryServiceClient) Update(ctx context.Context, id string, v *boundaries.PolicyBoundary) error {
	var err error

	client := iam.NewIAMClient(me)

	if _, err = client.PUT_MULTI_RESPONSE(
		ctx,
		fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id),
		v,
		[]int{201, 204},
		false,
	); err != nil {
		return err

	}

	return nil
}

func (me *BoundaryServiceClient) Delete(ctx context.Context, id string) error {
	var err error

	client := iam.NewIAMClient(me)

	if _, err = client.DELETE(
		ctx,
		fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id),
		204,
		false,
	); err != nil {
		if strings.Contains(err.Error(), "Policy boundary is in use") {
			clean.CleanUp.Register(func() {
				client.DELETE(
					ctx,
					fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, strings.TrimPrefix(me.AccountID(), "urn:dtaccount:"), id),
					204,
					false,
				)
			})
			return nil
		}
		return err
	}

	return nil
}

type PolicyBoundary struct {
	UUID      string `json:"uuid"`
	LevelType string `json:"levelType"`
	LevelID   string `json:"levelId"`
	Name      string `json:"name"`
	Query     string `json:"boundaryQuery"`
}

type ListPolicyBoundariesResponse struct {
	PageSize         int              `json:"pageSize"`
	PageNumber       int              `json:"pageNumber"`
	TotalCount       int              `json:"totalCount"`
	PolicyBoundaries []PolicyBoundary `json:"content"`
}
