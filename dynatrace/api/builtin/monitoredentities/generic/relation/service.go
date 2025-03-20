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

package relation

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	relation "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/monitoredentities/generic/relation/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "1.1"
const SchemaID = "builtin:monitoredentities.generic.relation"

func Service(credentials *rest.Credentials) settings.CRUDService[*relation.Settings] {
	return settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*relation.Settings]{Duplicates: Duplicates})
}

func Duplicates(ctx context.Context, service settings.RService[*relation.Settings], v *relation.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_generic_relationships") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*relation.Settings)
			if v.FromType == config.FromType && v.ToType == config.ToType {
				return nil, fmt.Errorf("Generic relationship with source/destination type already exists")
			}
		}
	} else if settings.HijackDuplicate("dynatrace_generic_relationships") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*relation.Settings)
			if v.FromType == config.FromType && v.ToType == config.ToType {
				return stub, nil
			}
		}
	}
	return nil, nil
}
