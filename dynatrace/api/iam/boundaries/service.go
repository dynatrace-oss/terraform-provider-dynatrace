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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	boundaries "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam/boundaries/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/clean"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const maxPageSize = 10000

func Service(credentials *rest.Credentials) settings.CRUDService[*boundaries.PolicyBoundary] {
	return &BoundaryServiceClient{
		iamClientGetter: &iamClientGetterImp{
			clientID:     credentials.IAM.ClientID,
			accountID:    credentials.IAM.AccountID,
			clientSecret: credentials.IAM.ClientSecret,
			tokenURL:     credentials.IAM.TokenURL,
			endpointURL:  credentials.IAM.EndpointURL,
		},
		accountID:   credentials.IAM.AccountID,
		endpointURL: credentials.IAM.EndpointURL,
	}
}

type iamClientGetter interface {
	New(ctx context.Context) iam.IAMClient
}

type iamClientGetterImp struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func (me *iamClientGetterImp) ClientID() string {
	return me.clientID
}

func (me *iamClientGetterImp) AccountID() string {
	return me.accountID
}

func (me *iamClientGetterImp) ClientSecret() string {
	return me.clientSecret
}

func (me *iamClientGetterImp) TokenURL() string {
	return me.tokenURL
}

func (me *iamClientGetterImp) EndpointURL() string {
	return me.endpointURL
}

func (me *iamClientGetterImp) New(ctx context.Context) iam.IAMClient {
	return iam.NewIAMClient(ctx, me)
}

type BoundaryServiceClient struct {
	iamClientGetter iamClientGetter
	accountID       string
	endpointURL     string
}

func (me *BoundaryServiceClient) List(ctx context.Context) (api.Stubs, error) {
	client := me.iamClientGetter.New(ctx)
	stubs := api.Stubs{}
	params := url.Values{}
	params.Set("size", strconv.Itoa(maxPageSize))

	for page := 1; ; page++ {
		params.Set("page", strconv.Itoa(page))
		reqURL := fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries?%s", me.endpointURL, me.accountID, params.Encode())
		responseBytes, err := client.GET(ctx, reqURL, 200, false)
		if err != nil {
			return nil, err
		}

		var response ListPolicyBoundariesResponse
		if err = json.Unmarshal(responseBytes, &response); err != nil {
			return nil, err
		}

		for _, boundary := range response.PolicyBoundaries {
			stubs = append(stubs, &api.Stub{ID: boundary.UUID, Name: boundary.Name})
		}

		if len(response.PolicyBoundaries) < maxPageSize {
			break
		}
	}

	return stubs, nil
}

func (me *BoundaryServiceClient) Get(ctx context.Context, id string, v *boundaries.PolicyBoundary) error {
	responseBytes, err := me.iamClientGetter.New(ctx).GET(ctx, fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, me.accountID, id), 200, false)
	if err != nil {
		return err
	}

	return json.Unmarshal(responseBytes, v)
}

func (me *BoundaryServiceClient) SchemaID() string {
	return "accounts:iam:boundaries"
}

func (me *BoundaryServiceClient) Create(ctx context.Context, v *boundaries.PolicyBoundary) (*api.Stub, error) {
	responseBytes, err := me.iamClientGetter.New(ctx).POST(
		ctx,
		fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries", me.endpointURL, me.accountID),
		v,
		201,
		false,
	)
	if err != nil {
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
	_, err := me.iamClientGetter.New(ctx).PUT_MULTI_RESPONSE(
		ctx,
		fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, me.accountID, id),
		v,
		[]int{201, 204},
		false,
	)
	return err
}

func (me *BoundaryServiceClient) Delete(ctx context.Context, id string) error {
	client := me.iamClientGetter.New(ctx)
	_, err := client.DELETE(
		ctx,
		fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, me.accountID, id),
		204,
		false,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Policy boundary is in use") {
			clean.CleanUp.Register(func() {
				client.DELETE(
					ctx,
					fmt.Sprintf("%s/iam/v1/repo/account/%s/boundaries/%s", me.endpointURL, me.accountID, id),
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
