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

package hosts

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	hosts "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/hosts/settings"
)

const SchemaID = "v1:config:conditional-naming:hosts"
const BasePath = "/api/config/v1/conditionalNaming/host"

func Service(credentials *settings.Credentials) settings.CRUDService[*hosts.NamingRule] {
	return settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*hosts.NamingRule](BasePath).WithDuplicates(Duplicates),
	)
}

func Duplicates(ctx context.Context, service settings.RService[*hosts.NamingRule], v *hosts.NamingRule) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_host_naming") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return nil, fmt.Errorf("A Host Naming Rule named '%s' already exists", v.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_host_naming") {
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
