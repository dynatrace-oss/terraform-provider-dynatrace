/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package boundaries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	boundaries "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/boundaries/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/clean"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

const maxPageSize = 10000

func Service(clientSet rest.ClientSet) (settings.CRUDService[*boundaries.PolicyBoundary], error) {
	return &BoundaryServiceClient{
		iamClientGetter: &iamClientGetterImp{
			credentials: clientSet.Credentials(),
		},
	}, nil
}

type iamClientGetter interface {
	New(ctx context.Context) rest.IAMClient
}

type iamClientGetterImp struct {
	credentials *rest.Credentials
}

func (me *iamClientGetterImp) New(ctx context.Context) rest.IAMClient {
	return rest.NewIAMClient(ctx, me.credentials)
}

type BoundaryServiceClient struct {
	iamClientGetter iamClientGetter
}

func (me *BoundaryServiceClient) List(ctx context.Context) (api.Stubs, error) {
	client := me.iamClientGetter.New(ctx)
	stubs := api.Stubs{}
	params := url.Values{}
	params.Set("size", strconv.Itoa(maxPageSize))
	endpoint := fmt.Sprintf("/iam/v1/repo/account/%s/boundaries", client.AccountID())

	for page := 1; ; page++ {
		params.Set("page", strconv.Itoa(page))
		response, err := client.GET(ctx, endpoint, rest2.RequestOptions{QueryParams: params})
		if err != nil {
			return nil, err
		}

		var policyBoundariesResponse ListPolicyBoundariesResponse
		if err = json.Unmarshal(response.Data, &policyBoundariesResponse); err != nil {
			return nil, err
		}

		for _, boundary := range policyBoundariesResponse.PolicyBoundaries {
			stubs = append(stubs, &api.Stub{ID: boundary.UUID, Name: boundary.Name})
		}

		if len(policyBoundariesResponse.PolicyBoundaries) < maxPageSize {
			break
		}
	}

	return stubs, nil
}

func (me *BoundaryServiceClient) Get(ctx context.Context, id string, v *boundaries.PolicyBoundary) error {
	client := me.iamClientGetter.New(ctx)
	response, err := client.GET(ctx, fmt.Sprintf("/iam/v1/repo/account/%s/boundaries/%s", client.AccountID(), id), rest2.RequestOptions{})
	if err != nil {
		return err
	}

	return json.Unmarshal(response.Data, v)
}

func (me *BoundaryServiceClient) SchemaID() string {
	return "accounts:iam:boundaries"
}

func (me *BoundaryServiceClient) Create(ctx context.Context, v *boundaries.PolicyBoundary) (*api.Stub, error) {
	client := me.iamClientGetter.New(ctx)
	response, err := client.POST(
		ctx,
		fmt.Sprintf("/iam/v1/repo/account/%s/boundaries", client.AccountID()),
		v,
		rest2.RequestOptions{},
	)
	if err != nil {
		return nil, err
	}

	uuidNameResponse := struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	}{}

	if err = json.Unmarshal(response.Data, &uuidNameResponse); err != nil {
		return nil, err
	}

	return &api.Stub{ID: uuidNameResponse.UUID, Name: uuidNameResponse.Name}, nil
}

func (me *BoundaryServiceClient) Update(ctx context.Context, id string, v *boundaries.PolicyBoundary) error {
	client := me.iamClientGetter.New(ctx)
	_, err := client.PUT(
		ctx,
		fmt.Sprintf("/iam/v1/repo/account/%s/boundaries/%s", client.AccountID(), id),
		v,
		rest2.RequestOptions{},
	)
	return err
}

func (me *BoundaryServiceClient) Delete(ctx context.Context, id string) error {
	client := me.iamClientGetter.New(ctx)
	_, err := client.DELETE(
		ctx,
		fmt.Sprintf("/iam/v1/repo/account/%s/boundaries/%s", client.AccountID(), id),
		rest2.RequestOptions{},
	)
	if err != nil {
		if strings.Contains(err.Error(), "Policy boundary is in use") {
			clean.CleanUp.Register(func() {
				client.DELETE(
					ctx,
					fmt.Sprintf("/iam/v1/repo/account/%s/boundaries/%s", client.AccountID(), id),
					rest2.RequestOptions{},
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
