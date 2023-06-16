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

package nodes

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	nodes "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/nodes/settings"
)

const SchemaID = "v1:synthetic:nodes:all"

func Service(credentials *settings.Credentials) settings.RService[*nodes.Settings] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

type service struct {
	client rest.Client
}

type nodeList struct {
	Nodes []nodes.Settings `json:"nodes"`
}

func (me *service) List() (stubs api.Stubs, err error) {
	var stubList nodeList
	if err = me.client.Get("/api/v1/synthetic/nodes", 200).Finish(&stubList); err != nil {
		return nil, err
	}
	for _, node := range stubList.Nodes {
		n := node
		stubs = append(stubs, &api.Stub{ID: node.ID, Name: node.Hostname, Value: n})
	}
	return stubs, nil
}

func (me *service) Get(id string, v *nodes.Settings) (err error) {
	return me.client.Get(fmt.Sprintf("/api/v1/synthetic/nodes/%v", id), 200).Finish(&v)
}

func (me *service) SchemaID() string {
	return SchemaID
}
