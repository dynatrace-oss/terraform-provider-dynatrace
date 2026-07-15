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
	"encoding/json"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/iam"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	rest2 "github.com/dynatrace/dynatrace-configuration-as-code-core/api/rest"
)

const baseURL = "/env/v2/accounts"

type environmentsService struct {
	credentials *rest.Credentials
}

func newEnvironmentService(clientSet rest.ClientSet) *environmentsService {
	return &environmentsService{
		credentials: clientSet.Credentials(),
	}
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
	u, err := url.JoinPath(baseURL, ec.credentials.IAM.AccountID, "environments")
	if err != nil {
		return nil, err
	}

	client := iam.NewIAMClient(ctx, ec.credentials)
	response, err := client.GET(ctx, u, rest2.RequestOptions{})
	if err != nil {
		return nil, err
	}

	var environmentResponse Response
	if err = json.Unmarshal(response.Data, &environmentResponse); err != nil {
		return nil, err
	}

	return environmentResponse.Data, nil
}
