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
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/azure/services/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:credentials:azure:services"
const BasePath = "/api/config/v1/azure/credentials"

var smu sync.Mutex

func Service(credentials *rest.Credentials) settings.CRUDService[*services.Settings] {
	return &service{
		client:     rest.APITokenClient(credentials),
		supService: NewSupportedServicesService(credentials),
	}
}

type service struct {
	client     rest.Client
	supService *SupportedServicesService
}

type srvStubs struct {
	Services []struct {
		Name string
	} `json:"services"`
}

type servicesResponse struct {
	Services []*services.Settings `json:"services"`
}

func cloneMetric(m *services.AzureMonitoredMetric) *services.AzureMonitoredMetric {
	cl := services.AzureMonitoredMetric{
		Name:       m.Name,
		Dimensions: []string{},
	}
	for _, d := range m.Dimensions {
		cl.Dimensions = append(cl.Dimensions, d)
	}
	return &cl
}

func cloneMetrics(m []*services.AzureMonitoredMetric) []*services.AzureMonitoredMetric {
	cl := []*services.AzureMonitoredMetric{}
	for _, m := range m {
		cl = append(cl, cloneMetric(m))
	}
	return cl
}

func cloneService(s *services.Settings) *services.Settings {
	return &services.Settings{
		CredentialsID:    s.CredentialsID,
		ServiceName:      s.ServiceName,
		MonitoredMetrics: cloneMetrics(s.MonitoredMetrics),
		BuiltIn:          s.BuiltIn,
		RequiredMetrics:  s.RequiredMetrics,
	}
}

func (me servicesResponse) clone() servicesResponse {
	var cl servicesResponse
	for _, s := range me.Services {
		cl.Services = append(cl.Services, cloneService(s))
	}
	return cl
}

type listResponse struct {
	Values api.Stubs `json:"values"`
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs
	var credentialStubs listResponse
	var err error
	if err = me.client.Get(ctx, "/api/config/v1/azure/credentials").Expect(200).Finish(&credentialStubs); err != nil {
		return nil, err
	}
	for _, credentialStub := range credentialStubs.Values {
		var servicesStubs srvStubs
		if err = me.client.Get(ctx, fmt.Sprintf("/api/config/v1/azure/credentials/%s/services", credentialStub.ID)).Expect(200).Finish(&servicesStubs); err != nil {
			return nil, err
		}
		for _, servicesStub := range servicesStubs.Services {
			stubs = append(stubs, &api.Stub{ID: credentialStub.ID + "#" + servicesStub.Name, Name: credentialStub.ID + "_" + servicesStub.Name})
		}
	}
	return stubs, nil
}

func (me *service) Get(ctx context.Context, id string, v *services.Settings) error {
	smu.Lock()
	defer smu.Unlock()
	parts := strings.Split(id, "#")
	credentialsID := parts[0]
	serviceName := parts[1]
	var response servicesResponse
	var err error
	if err = me.client.Get(ctx, fmt.Sprintf("/api/config/v1/azure/credentials/%s/services", credentialsID)).Expect(200).Finish(&response); err != nil {
		return err
	}
	for _, service := range response.Services {
		if strings.ToLower(service.ServiceName) == strings.ToLower(serviceName) {
			v.CredentialsID = credentialsID
			v.MonitoredMetrics = service.MonitoredMetrics
			v.ServiceName = serviceName
			v.BuiltIn, _ = me.supService.IsBuiltIn(ctx, v.ServiceName)
			return nil
		}
	}
	return rest.Error{Code: 404, Message: fmt.Sprintf("No service '%s' for credential '%s' found", serviceName, credentialsID)}
}

func (me *service) SchemaID() string {
	return SchemaID
}

func (me *service) Create(ctx context.Context, v *services.Settings) (*api.Stub, error) {
	smu.Lock()
	defer smu.Unlock()
	credentialsID := v.CredentialsID
	var response servicesResponse
	var err error
	if err = me.client.Get(ctx, fmt.Sprintf("/api/config/v1/azure/credentials/%s/services", credentialsID)).Expect(200).Finish(&response); err != nil {
		return nil, err
	}
	if v.UseRecommendedMetrics {
		v.MonitoredMetrics = nil
	}

	isBuiltIn, e := me.supService.IsBuiltIn(ctx, strings.ToLower(v.ServiceName))
	if e != nil {
		return nil, e
	}
	var service *services.Settings
	for _, s := range response.Services {
		if strings.ToLower(s.ServiceName) == strings.ToLower(v.ServiceName) {
			service = s
			break
		}
	}
	if service != nil {
		if isBuiltIn || v.UseRecommendedMetrics {
			service.MonitoredMetrics = nil
		} else {
			service.MonitoredMetrics = v.MonitoredMetrics
		}
	} else {
		if isBuiltIn {
			v.MonitoredMetrics = nil
		}
		response.Services = append(response.Services, v)
	}

	retry := true
	for retry {
		if err = me.client.Put(ctx, fmt.Sprintf("/api/config/v1/azure/credentials/%s/services", credentialsID), response).Expect(204).Finish(); err != nil {
			r := regexp.MustCompile("Invalid\\sservices\\sconfiguration\\:\\srecommended\\smetrics\\s\\[([^\\]]*)\\]\\sfor\\sservice\\s'([^']*)'\\smust\\sbe\\sselected")
			r2 := regexp.MustCompile("Invalid\\sservices\\sconfiguration\\:\\smetric\\s'([^']*)'\\sfor\\sservice\\s'([^']*)'\\shas\\smissing\\sdimension\\s\\[([^\\]]*)\\],\\suse\\sall\\srecommended\\sdimensions\\s\\[([^\\]]*)\\]")
			if m := r.FindStringSubmatch(err.Error()); m != nil {
				var service *services.Settings
				for _, service = range response.Services {
					if strings.ToLower(service.ServiceName) == strings.ToLower(v.ServiceName) {
						break
					}
				}
				for _, metric := range strings.Split(strings.TrimSpace(m[1]), ",") {
					metricName := strings.TrimSpace(metric)
					service.MonitoredMetrics = append(service.MonitoredMetrics, &services.AzureMonitoredMetric{Name: metricName})
					if len(v.RequiredMetrics) == 0 {
						v.RequiredMetrics = metricName
					} else {
						v.RequiredMetrics = v.RequiredMetrics + "," + metricName
					}
				}
				v.MonitoredMetrics = service.MonitoredMetrics
			} else if m := r2.FindStringSubmatch(err.Error()); m != nil {
				var service *services.Settings
				for _, service = range response.Services {
					if strings.ToLower(service.ServiceName) == strings.ToLower(v.ServiceName) {
						break
					}
				}
				metricName := m[1]
				sDimensions := m[4]
				var metric *services.AzureMonitoredMetric
				for _, met := range service.MonitoredMetrics {
					if met.Name == metricName {
						metric = met
						break
					}
				}
				dimMap := map[string]string{}
				for _, k := range metric.Dimensions {
					dimMap[k] = k
				}
				for _, d := range strings.Split(sDimensions, ",") {
					td := strings.TrimSpace(d)
					dimMap[td] = td
				}
				metric.Dimensions = []string{}
				for k := range dimMap {
					metric.Dimensions = append(metric.Dimensions, k)
				}
				v.MonitoredMetrics = service.MonitoredMetrics
			} else {
				return nil, err
			}
		} else {
			retry = false
		}
	}

	return &api.Stub{ID: credentialsID + "#" + v.ServiceName, Name: credentialsID + "_" + v.ServiceName}, nil
}

func (me *service) Update(ctx context.Context, id string, v *services.Settings) error {
	_, err := me.Create(ctx, v)
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	smu.Lock()
	defer smu.Unlock()
	parts := strings.Split(id, "#")
	credentialsID := parts[0]
	serviceName := parts[1]
	var response servicesResponse
	var err error
	if err = me.client.Get(ctx, fmt.Sprintf("/api/config/v1/azure/credentials/%s/services", credentialsID)).Expect(200).Finish(&response); err != nil {
		return err
	}
	var reducedServices servicesResponse
	found := false
	for _, service := range response.Services {
		if strings.ToLower(service.ServiceName) == strings.ToLower(serviceName) {
			found = true
		} else {
			reducedServices.Services = append(reducedServices.Services, service)
		}
	}
	if !found {
		return nil
	}
	if len(reducedServices.Services) == 0 {
		reducedServices.Services = []*services.Settings{}
	}
	return me.client.Put(ctx, fmt.Sprintf("/api/config/v1/azure/credentials/%s/services", credentialsID), reducedServices).Expect(204).Finish()
}
