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

package services

import (
	"context"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/httpcache"
)

const ServiceSchemaID = "v1:config:credentials:aws:services:all"

var mu sync.Mutex

type SupportedService struct {
	CloudProviderServiceType string `json:"cloudProviderServiceType"`
	Name                     string `json:"name"`
	EntityType               string `json:"entityType"`
	DisplayName              string `json:"displayName"` // suffix `(built-in)` for built in services
	BuiltIn                  bool   `json:"-"`
}

var supportedServicesRepo = map[string]map[string]*SupportedService{}

func NewSupportedServicesService(credentials *settings.Credentials) *SupportedServicesService {
	return &SupportedServicesService{
		url:    credentials.URL,
		client: httpcache.DefaultClient(credentials.URL, credentials.Token, SchemaID),
	}
}

type SupportedServicesService struct {
	url    string
	client rest.Client
}

type supportedServicesResponse struct {
	Services []*SupportedService `json:"services"`
}

func (me *SupportedServicesService) IsBuiltIn(name string) (bool, error) {
	s, e := me.Get(context.Background(), name)
	if e != nil {
		return false, e
	}
	return s != nil && s.BuiltIn, nil
}

func (me *SupportedServicesService) Get(ctx context.Context, name string) (*SupportedService, error) {
	services, err := me.List(ctx)
	if err != nil {
		return nil, err
	}
	result, found := services[strings.ToLower(name)]
	if !found {
		return nil, nil
	}
	return result, nil
}

func (me *SupportedServicesService) List(ctx context.Context) (map[string]*SupportedService, error) {
	mu.Lock()
	defer mu.Unlock()

	if stored, found := supportedServicesRepo[me.url]; found {
		return stored, nil
	}
	var servicesResponse supportedServicesResponse
	if err := me.client.Get("/api/config/v1/aws/supportedServices").Expect(200).Finish(&servicesResponse); err != nil {
		return nil, err
	}
	if servicesResponse.Services == nil {
		servicesResponse.Services = []*SupportedService{}
	} else {
		for _, entry := range servicesResponse.Services {
			entry.BuiltIn = strings.HasSuffix(entry.DisplayName, "(built-in)")
			if entry.BuiltIn {
				entry.Name = strings.ToLower(entry.Name)
			}
		}
	}
	result := map[string]*SupportedService{}
	for _, entry := range servicesResponse.Services {
		result[strings.ToLower(entry.Name)] = entry
	}
	supportedServicesRepo[me.url] = result
	return result, nil
}
