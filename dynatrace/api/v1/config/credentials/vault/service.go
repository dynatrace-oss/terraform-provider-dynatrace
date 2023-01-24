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

package vault

import (
	"fmt"
	"net/url"
	"strings"

	vault "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/vault/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:credentials"
const BasePath = "/api/config/v1/credentials"

func Service(credentials *settings.Credentials) settings.CRUDService[*vault.Credentials] {
	return settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*vault.Credentials](BasePath).
			WithStubs(&vault.CredentialsList{}).
			NoValidator().
			WithDeleteRetry(func(id string, err error) (bool, error) {
				if strings.Contains(err.Error(), "as long as there are monitors assigned to it") {
					client := rest.DefaultClient(credentials.URL, credentials.Token)
					response := struct {
						Monitors []struct {
							EntityID string `json:"entityId"`
						} `json:"monitors"`
					}{}
					if err := client.Get(fmt.Sprintf("/api/v1/synthetic/monitors?credentialId=%s", url.QueryEscape(id)), 200).Finish(&response); err != nil {
						return false, err
					}
					if len(response.Monitors) > 0 {
						for _, monitor := range response.Monitors {
							if err := client.Delete(fmt.Sprintf("/api/v1/synthetic/monitors/%s", url.PathEscape(monitor.EntityID)), 204).Finish(); err != nil {
								return false, err
							}
						}
					}
					return true, nil
				}
				return false, nil
			}),
	)
}
