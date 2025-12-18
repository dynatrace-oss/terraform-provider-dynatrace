/*
 * @license
 * Copyright 2025 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package environments

import (
	"context"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
)

const baseURL = "/env/v2/accounts"

type environmentsService struct {
	clientID     string
	accountID    string
	clientSecret string
	tokenURL     string
	endpointURL  string
}

func newEnvironmentService(credentials *rest.Credentials) *environmentsService {
	return &environmentsService{
		clientID:     credentials.IAM.ClientID,
		accountID:    credentials.IAM.AccountID,
		clientSecret: credentials.IAM.ClientSecret,
		tokenURL:     credentials.IAM.TokenURL,
		endpointURL:  credentials.IAM.EndpointURL,
	}
}

func (ec *environmentsService) ClientID() string {
	return ec.clientID
}

func (ec *environmentsService) AccountID() string {
	return ec.accountID
}

func (ec *environmentsService) ClientSecret() string {
	return ec.clientSecret
}

func (ec *environmentsService) TokenURL() string {
	return ec.tokenURL
}

func (ec *environmentsService) EndpointURL() string {
	return ec.endpointURL
}

type Environment struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	URL    string `json:"url"`
}

type Response struct {
	Data []Environment `json:"data"`
}

func (ec *environmentsService) Get(ctx context.Context) ([]Environment, error) {
	u, err := url.JoinPath(ec.endpointURL, baseURL, ec.AccountID(), "environments")

	if err != nil {
		return nil, err
	}

	client := iam.NewIAMClient(ec)
	var result Response
	err = iam.GET(client, ctx, u, 200, false, &result)
	if err != nil {
		return nil, err
	}

	return result.Data, nil
}
