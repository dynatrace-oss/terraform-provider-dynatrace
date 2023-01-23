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

package customservices

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	customservices "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customservices/settings"
)

const SchemaID = "v1:config:custom-services"

type service struct {
	client rest.Client
}

func Service(credentials *settings.Credentials) settings.CRUDService[*customservices.CustomService] {
	return &service{client: rest.DefaultClient(credentials.URL, credentials.Token)}
}

func (me *service) Get(id string, v *customservices.CustomService) error {
	if id, technology, ok := settings.SplitID(id); ok {
		return me.GetWithTechnology(id, technology, v)
	}
	return errors.New("unable to determine technology")
}

func (me *service) GetWithTechnology(id string, technology string, v *customservices.CustomService) error {
	req := me.client.Get(fmt.Sprintf("/api/config/v1/service/customServices/%s/%s", url.PathEscape(technology), url.PathEscape(id))).Expect(200)
	if err := req.Finish(v); err != nil {
		return err
	}
	v.Technology = customservices.Technology(technology)
	return nil
}

func (me *service) List() (settings.Stubs, error) {
	var err error
	var stubs settings.Stubs
	client := me.client

	for _, technology := range []customservices.Technology{customservices.Technologies.NodeJS, customservices.Technologies.DotNet, customservices.Technologies.Go, customservices.Technologies.Java, customservices.Technologies.PHP} {
		req := client.Get(fmt.Sprintf("/api/config/v1/service/customServices/%s", url.PathEscape(string(technology))), 200)
		var stubList settings.StubList
		if err = req.Finish(&stubList); err != nil {
			return nil, err
		}
		for _, stub := range stubList.Values {
			stub.ID = settings.JoinID(stub.ID, string(technology))
			stubs = append(stubs, stub)
		}
	}
	return stubs, nil
}

func (me *service) Validate(v *customservices.CustomService) error {
	return me.ValidateWithTechnology(string(v.Technology), v)
}

func (me *service) ValidateWithTechnology(technology string, v any) error {
	return me.client.Post(fmt.Sprintf("/api/config/v1/service/customServices/%s/validator", url.PathEscape(technology)), v, 204).Finish()
}

func (me *service) Create(v *customservices.CustomService) (*settings.Stub, error) {
	return me.CreateWithTechnology(string(v.Technology), v)
}

func (me *service) CreateWithTechnology(technology string, v any) (*settings.Stub, error) {
	var err error

	req := me.client.Post(fmt.Sprintf("/api/config/v1/service/customServices/%s", url.PathEscape(technology)), v, 201)

	var stub settings.Stub
	if err = req.Finish(&stub); err != nil {
		return nil, err
	}
	stub.ID = settings.JoinID(stub.ID, technology)
	return &stub, nil
}

func (me *service) Update(id string, v *customservices.CustomService) error {
	if id, technology, ok := settings.SplitID(id); ok {
		return me.UpdateWithTechnology(id, technology, v)
	}
	return me.UpdateWithTechnology(id, string(v.Technology), v)
}

func (me *service) UpdateWithTechnology(id string, technology string, v any) error {
	var err error

	req := me.client.Put(fmt.Sprintf("/api/config/v1/service/customServices/%s/%s", url.PathEscape(technology), url.PathEscape(id)), v, 204)

	if err = req.Finish(); err != nil {
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	if id, technology, ok := settings.SplitID(id); ok {
		return me.DeleteWithTechnology(id, technology)
	}
	return errors.New("unable to determine technology")
}

func (me *service) DeleteWithTechnology(id string, technology string) error {
	return me.client.Delete(fmt.Sprintf("/api/config/v1/service/customServices/%s/%s", url.PathEscape(technology), url.PathEscape(id))).Expect(204).Finish()
}

func (me *service) SchemaID() string {
	return SchemaID
}
