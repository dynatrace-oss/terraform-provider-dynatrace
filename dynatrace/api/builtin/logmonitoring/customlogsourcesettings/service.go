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

package customlogsourcesettings

import (
	"context"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	customlogsourcesettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/logmonitoring/customlogsourcesettings/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.1"
const SchemaID = "builtin:logmonitoring.custom-log-source-settings"

// The custom service below was created to gracefully deal with the breaking change from attribute `values` to `values-and-enrichment` introduced in v292.
func Service(credentials *settings.Credentials) settings.CRUDService[*customlogsourcesettings.Settings] {
	return &service{
		service: settings20.Service[*customlogsourcesettings.Settings](credentials, SchemaID, SchemaVersion),
		client:  rest.DefaultClient(credentials.URL, credentials.Token),
	}
}

type service struct {
	service settings.CRUDService[*customlogsourcesettings.Settings]
	client  rest.Client
}

var errorMessage = "Given property 'values-and-enrichment' with value: '0' violates the following constraint: fell below the collection's lower size limit which was set to 1."

func (me *service) Create(v *customlogsourcesettings.Settings) (*api.Stub, error) {
	var stub *api.Stub
	var err error

	if stub, err = me.service.Create(v); err != nil {
		if strings.Contains(err.Error(), errorMessage) && v.Custom_log_source != nil && len(v.Custom_log_source.Values) > 0 {
			valuesToValuesAndEnrichment(v)
			if stub, err = me.service.Create(v); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return stub, nil
}

func (me *service) Update(id string, v *customlogsourcesettings.Settings) error {
	if err := me.service.Update(id, v); err != nil {
		if strings.Contains(err.Error(), errorMessage) && v.Custom_log_source != nil && len(v.Custom_log_source.Values) > 0 {
			valuesToValuesAndEnrichment(v)
			if err = me.service.Update(id, v); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (me *service) GetWithContext(ctx context.Context, id string, v *customlogsourcesettings.Settings) error {
	var err error
	if err = me.Get(id, v); err != nil {
		return err
	}

	if v.Custom_log_source != nil && len(v.Custom_log_source.Values) == 0 {
		var values []string
		cfg := ctx.Value(settings.ContextKeyStateConfig)
		if logConfig, ok := cfg.(*customlogsourcesettings.Settings); ok {
			values = logConfig.Custom_log_source.Values
		}

		if len(values) > 0 {
			valuesAndEnrichmentToValues(v)
		}
	}

	return nil
}

func valuesToValuesAndEnrichment(v *customlogsourcesettings.Settings) {
	for _, value := range v.Custom_log_source.Values {
		temp := new(customlogsourcesettings.CustomLogSourceWithEnrichment)
		temp.Path = value
		v.Custom_log_source.Values_and_enrichment = append(v.Custom_log_source.Values_and_enrichment, temp)
	}
	v.Custom_log_source.Values = []string{}
}

func valuesAndEnrichmentToValues(v *customlogsourcesettings.Settings) {
	for _, value := range v.Custom_log_source.Values_and_enrichment {
		v.Custom_log_source.Values = append(v.Custom_log_source.Values, value.Path)
	}
	v.Custom_log_source.Values_and_enrichment = customlogsourcesettings.CustomLogSourceWithEnrichments{}
}

func (me *service) Get(id string, v *customlogsourcesettings.Settings) error {
	return me.service.Get(id, v)
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func (me *service) List() (api.Stubs, error) {
	return me.service.List()
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}
