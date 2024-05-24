/**
* @license
* Copyright 2020 Dynatrace LLC
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

package lambdaagent

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	lambdaagent "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/deployment/lambdaagent/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:deployment:lambdaagent"
const BasePath = "/api/v1/deployment/lambda/agent/latest"

func Service(credentials *settings.Credentials) settings.RService[*lambdaagent.Latest] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

func (me *service) Get(id string, v *lambdaagent.Latest) error {
	if err := me.client.Get(BasePath, 200).Finish(&v); err != nil {
		return err
	}
	return nil
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) List() (api.Stubs, error) {
	return api.Stubs{&api.Stub{ID: me.SchemaID(), Name: me.SchemaID()}}, nil
}
