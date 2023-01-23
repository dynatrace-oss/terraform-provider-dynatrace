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

package networkzones_test

import (
	"os"
	"testing"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/networkzones"
	v2n "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/networkzones"
	v2networkzones "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/networkzones/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing/api"
)

func TestNetworkZones(t *testing.T) {
	var err error
	var stub *settings.Stub
	var cfg v2networkzones.NetworkZones
	var stubs settings.Stubs
	envURL := os.Getenv("DYNATRACE_ENV_URL")
	apiToken := os.Getenv("DYNATRACE_API_TOKEN")
	if envURL == "" || apiToken == "" {
		t.Skip("Environment Variables DYNATRACE_ENV_URL and DYNATRACE_API_TOKEN must be specified")
		return
	}

	service := v2n.Service(&settings.Credentials{URL: envURL, Token: apiToken})
	if stubs, err = service.List(); err != nil {
		t.Error(err)
		return
	}
	if len(stubs) == 0 {
		if stub, err = service.Create(&v2networkzones.NetworkZones{Enabled: true}); err != nil {
			t.Error(err)
			return
		}
	} else {
		if err = service.Get(stubs[0].ID, &cfg); err != nil {
			t.Error(err)
			return
		}
		if err = service.Update(stubs[0].ID, &v2networkzones.NetworkZones{Enabled: true}); err != nil {
			t.Error(err)
			return
		}
	}

	for {
		var sts settings.Stubs
		var scfg v2networkzones.NetworkZones
		if sts, err = service.List(); err != nil {
			t.Error(err)
			return
		}
		if len(sts) > 0 {
			if err = service.Get(sts[0].ID, &scfg); err != nil {
				t.Error(err)
				return
			}
			if scfg.Enabled {
				break
			}
		}
	}

	time.Sleep(5 * time.Second)

	api.TestService(t, networkzones.Service)
	if stub != nil {
		if err = service.Delete(stub.ID); err != nil {
			t.Error(err)
			return
		}
	} else {
		if err = service.Update(stubs[0].ID, &cfg); err != nil {
			t.Error(err)
			return
		}
	}
}
