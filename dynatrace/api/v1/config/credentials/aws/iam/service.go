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

package iam

import (
	iam "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/aws/iam/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:credentials:aws:iam:external:id"
const StaticID = "7085923b-6268-441f-98df-90caf9aab7e5"
const StaticName = "aws_iam_external_id"

type service struct {
	client rest.Client
}

func (me *service) Get(id string, v *iam.Settings) error {
	return me.client.Get("/api/config/v1/aws/iamExternalId", 200).Finish(v)
}

func (me *service) List() (settings.Stubs, error) {
	return settings.Stubs{&settings.Stub{ID: StaticID, Name: StaticName}}, nil
}

func (me *service) SchemaID() string {
	return SchemaID
}

func Service(credentials *settings.Credentials) settings.RService[*iam.Settings] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}
