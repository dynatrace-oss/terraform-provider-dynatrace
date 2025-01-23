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

package notifications

import (
	"context"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/filtered"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaID = "builtin:problem.notifications"
const SchemaVersion = "1.6.1"

type filter struct {
	Type Type
}

func (me *filter) Filter(v *Notification) (bool, error) {
	return (string(v.Type) == string(me.Type)), nil
}

func (me *filter) Suffix() string {
	return string(me.Type)
}

func Service(credentials *settings.Credentials, t Type) settings.CRUDService[*Notification] {
	return filtered.Service[*Notification](
		settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*Notification]{
			LegacyID:    settings.LegacyObjIDDecode,
			CreateRetry: UseOAuth2Fix,
			UpdateRetry: UseOAuth2Fix,
			Duplicates:  Duplicates,
		}),
		&filter{Type: t},
	)
}

func UseOAuth2Fix(v *Notification, err error) *Notification {
	if strings.Contains(err.Error(), "Given property 'useOAuth2' with value: 'null' does not comply with required NonNull of schema") {
		v.WebHook.UseOAuth2 = opt.NewBool(false)
		return v
	}
	return nil
}

func Duplicates(ctx context.Context, service settings.RService[*Notification], v *Notification) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_notification") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(ctx); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			if v.Name == stub.Name {
				return nil, fmt.Errorf("Notification named '%s' already exists", v.Name)
			}
		}
	} else if settings.HijackDuplicate("dynatrace_notification") {
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
