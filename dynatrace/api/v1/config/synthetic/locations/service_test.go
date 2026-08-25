//go:build integration

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
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations"
	locsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/locations/settings"
	testing2 "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/testing"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
)

func TestSyntheticLocations(t *testing.T) {
	pc := config.ProviderConfigureGeneric(t.Context(), testing2.EnvironmentGetter(t))
	service, err := locations.Service(pc)
	if err != nil {
		t.Error(err)
		return
	}

	var stubs api.Stubs
	if stubs, err = service.List(context.Background()); err != nil {
		t.Error(err)
		return
	}

	foundPublic := false
	foundPrivate := false
	for _, stub := range stubs {
		if stub.Value != nil {
			loc := stub.Value.(*locsettings.SyntheticLocation)
			if loc.Type == locsettings.LocationTypes.Public {
				foundPublic = true
			} else if loc.Type == locsettings.LocationTypes.Private {
				foundPrivate = true
			}
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
