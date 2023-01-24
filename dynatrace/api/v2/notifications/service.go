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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/filtered"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaID = "builtin:problem.notifications"
const SchemaVersion = "1.4.1"

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
		settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*Notification]{LegacyID: settings.LegacyObjIDDecode}),
		&filter{Type: t},
	)
}
