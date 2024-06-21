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

package entrypoints

import (
	"context"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"

	entrypoints "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/span/entrypoints/settings"
)

const SchemaID = "builtin:span-entry-points"
const SchemaVersion = "0.1.16"

func Service(credentials *settings.Credentials) settings.CRUDService[*entrypoints.SpanEntryPoint] {
	return settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*entrypoints.SpanEntryPoint]{Duplicates: Duplicates})
}

func Duplicates(ctx context.Context, service settings.RService[*entrypoints.SpanEntryPoint], v *entrypoints.SpanEntryPoint) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_span_entry_point") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.EntryPointRule.Name == stub.Name {
				return nil, fmt.Errorf("An entry point rule named '%s' already exists", v.EntryPointRule.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_span_entry_point") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.EntryPointRule.Name == stub.Name {
				return stub, nil
			}
		}
	}
	return nil, nil
}
