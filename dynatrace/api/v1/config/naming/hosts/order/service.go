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

package hostnaming_order

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	order "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/order"
	order_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/naming/order/settings"
)

const SchemaID = "v1:config:conditional-naming:hosts:order"
const StaticID = "1551518f-4591-4c64-91ff-57f4b14dbb2e"
const StaticName = "host_naming_order"
const BasePath = "/api/config/v1/conditionalNaming/host"

func Service(credentials *rest.Credentials) settings.CRUDService[*order_settings.Settings] {
	return order.Service(credentials, SchemaID, StaticID, StaticName, BasePath)
}
