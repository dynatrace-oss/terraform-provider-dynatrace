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

package locations_test

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations"
	locsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

func TestSyntheticLocations(t *testing.T) {
	envURL := os.Getenv("DYNATRACE_ENV_URL")
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	if envURL == "" || apiToken == "" {
		t.Skip("Environment Variables DYNATRACE_ENV_URL and DYNATRACE_API_TOKEN must be specified")
		return
	}
	credentials := &settings.Credentials{URL: envURL, Token: apiToken}
	service := locations.Service(credentials)
	var stubs settings.Stubs
	var err error
	if stubs, err = service.List(); err != nil {
		t.Error(err)
		return
	}
	foundPublic := false
	foundPrivate := false
	for _, stub := range stubs {
		if strings.HasPrefix(stub.ID, "GEOLOCATION-") {
			foundPublic = true
		} else if strings.HasPrefix(stub.ID, "SYNTHETIC_LOCATION-") {
			foundPrivate = true
		}
		if stub.Value == nil {
			t.Error("Stubs were expected to contain values, but didn't")
			return
		}
		if stub.Value.(*locsettings.SyntheticLocation).ID != stub.ID {
			data, _ := json.Marshal(stub)
			log.Println("stub: " + string(data))
			data, _ = json.Marshal(stub.Value)
			log.Println("value: " + string(data))
			t.Error("ID of Stubs don't match ID of values")
			return
		}
	}
	if !foundPublic {
		t.Error("Expected to find public synthetic locations - found none")
		return
	}
	if !foundPrivate {
		t.Error("Expected to find private synthetic locations - found none")
		return
	}
}
