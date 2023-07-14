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

package aws

import (
	"fmt"
	"sync"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	aws "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/aws/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

const SchemaID = "v1:config:credentials:aws"
const BasePath = "/api/config/v1/aws/credentials"

var mu sync.Mutex

func Service(credentials *settings.Credentials) settings.CRUDService[*aws.AWSCredentialsConfig] {
	return settings.NewCRUDService(
		credentials,
		SchemaID,
		settings.DefaultServiceOptions[*aws.AWSCredentialsConfig](BasePath).
			WithStubs(&api.Stubs{}).
			WithMutex(mu.Lock, mu.Unlock).
			WithOnBeforeUpdate(func(id string, v *aws.AWSCredentialsConfig) error {
				if v.SupportingServicesManagedInDynatrace {
					var creds aws.AWSCredentialsConfig
					client := rest.DefaultClient(credentials.URL, credentials.Token)

					if err := client.Get(fmt.Sprintf("%s/%s", BasePath, "id"), 200).Finish(&creds); err != nil {
						return err
					}
					v.SupportingServicesToMonitor = creds.SupportingServicesToMonitor
				}
				return nil
			}).
			WithAfterCreate(func(client rest.Client, stub *api.Stub) (*api.Stub, error) {
				// After creating AWS Credentials it may take a while until the `externalId` has been set by the cluster
				// We're polling roughly 60 seconds until that has happened - in order to ensure that the credentials REALLY have been created
				// Last resort is to populate that value right after GET has happened (see below) - which is already cheating
				var cfg aws.AWSCredentialsConfig
				numRetries := 0
				configIsValid := false
				for !configIsValid && numRetries < 30 {
					client.Get(fmt.Sprintf("/api/config/v1/aws/credentials/%s", stub.ID), 200).Finish(&cfg)
					if cfg.AuthenticationData == nil || cfg.AuthenticationData.RoleBasedAuthentication == nil || (cfg.AuthenticationData.RoleBasedAuthentication.ExternalID != nil && len(*cfg.AuthenticationData.RoleBasedAuthentication.ExternalID) > 0) {
						configIsValid = true
						break
					}
					numRetries++
					time.Sleep(2 * time.Second)
				}
				return stub, nil
			}).
			WithCompleteGet(func(client rest.Client, id string, v *aws.AWSCredentialsConfig) error {
				// This is a sanity (last resort) function
				// Sometimes freshly created AWS Credentials don't have the `externalId` assigned yet
				// ... even after a miniute of waiting
				// If all of that fails, we're correcting that right after we've fetched the current state
				// because that `externalId` IS available globally elsewhere
				if v.AuthenticationData == nil {
					return nil
				}
				if v.AuthenticationData.RoleBasedAuthentication == nil {
					return nil
				}
				if (v.AuthenticationData.RoleBasedAuthentication.ExternalID == nil) || len(*v.AuthenticationData.RoleBasedAuthentication.ExternalID) == 0 {
					tokenResponse := struct {
						Token string `json:"token"`
					}{}
					client.Get("/api/config/v1/aws/iamExternalId", 200).Finish(&tokenResponse)
					if len(tokenResponse.Token) > 0 {
						v.AuthenticationData.RoleBasedAuthentication.ExternalID = &tokenResponse.Token
					}
				}
				return nil
			}),
	)
}
