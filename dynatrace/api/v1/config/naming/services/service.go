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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	services "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/services/settings"
)

const SchemaID = "v1:config:conditional-naming:services"
const BasePath = "/api/config/v1/conditionalNaming/service"

func Service(credentials *rest.Credentials) settings.CRUDService[*services.NamingRule] {
	return settings.NewAPITokenService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*services.NamingRule](BasePath).WithDuplicates(Duplicates),
	)
}

func Duplicates(ctx context.Context, service settings.RService[*services.NamingRule], v *services.NamingRule) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_service_naming") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return nil, fmt.Errorf("A Service Naming Rule named '%s' already exists", v.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_service_naming") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return stub, nil
			}
		}
	}
	return nil, nil
}
