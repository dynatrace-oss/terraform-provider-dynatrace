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

package azure

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	azure "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/azure/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
)

const SchemaID = "v1:config:credentials:azure"
const BasePath = "/api/config/v1/azure/credentials"

var mu sync.Mutex

func Service(credentials *settings.Credentials) settings.CRUDService[*azure.AzureCredentials] {
	return &service{service: settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*azure.AzureCredentials](BasePath).
			WithMutex(mu.Lock, mu.Unlock),
	), client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID)}
}

type service struct {
	service settings.CRUDService[*azure.AzureCredentials]
	client  rest.Client
}

func (me *service) List() (api.Stubs, error) {
	return me.service.List()
}

func (me *service) Get(id string, v *azure.AzureCredentials) error {
	return me.service.Get(id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(v *azure.AzureCredentials) (*api.Stub, error) {
	return me.service.Create(v)
}

var updateRegex = regexp.MustCompile("Invalid generic services configuration: You can't set up (.*) services using this endpoint")

type blackList []string

func (me blackList) Allows(v *azure.AzureSupportingService) bool {
	if len(me) == 0 {
		return true
	}
	for _, elem := range me {
		if v.Name != nil && strings.TrimSpace(elem) == strings.TrimSpace(*v.Name) {
			return false
		}
	}
	return true
}

func (me *service) Update(id string, v *azure.AzureCredentials) error {
	if v.SupportingServicesManagedInDynatrace {
		var creds azure.AzureCredentials
		if err := me.client.Get(fmt.Sprintf("%s/%s", BasePath, id), 200).Finish(&creds); err != nil {
			return err
		}
		v.SupportingServices = creds.SupportingServices
	}
	if err := me.service.Update(id, v); err != nil {
		if matches := updateRegex.FindStringSubmatch(err.Error()); len(matches) > 1 {
			if len(v.SupportingServices) == 0 {
				return err
			}
			var serviceStr string
			if serviceStr = matches[1]; len(serviceStr) == 0 {
				return err
			}
			blacklist := blackList{}
			for _, service := range strings.Split(serviceStr, ",") {
				blacklist = append(blacklist, strings.TrimSpace(service))
			}
			supportingServices := []*azure.AzureSupportingService{}
			for _, supportingservice := range v.SupportingServices {
				if blacklist.Allows(supportingservice) {
					supportingServices = append(supportingServices, supportingservice)
				}
			}
			if len(v.SupportingServices) == len(supportingServices) {
				return err
			}
			v.SupportingServices = supportingServices
			return me.service.Update(id, v)
		}
		return err
	}
	return nil
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) Name() string {
	return me.service.Name()
}
