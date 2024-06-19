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

package appdetection

import (
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	appdetection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/rum/web/appdetection/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "2.1.1"
const SchemaID = "builtin:rum.web.app-detection"

func Service(credentials *settings.Credentials) settings.CRUDService[*appdetection.Settings] {
	return &service{
		service: settings20.Service(credentials, SchemaID, SchemaVersion, &settings20.ServiceOptions[*appdetection.Settings]{Duplicates: Duplicates}),
		client:  rest.DefaultClient(credentials.URL, credentials.Token),
	}
}

var NotAValidDataSourceRegexp = regexp.MustCompile(`Given property 'applicationId' with value: '([^']*)' is not a valid value in datasource\. Value must be one of`)
var ApplicationIDRegexp = regexp.MustCompile("^APPLICATION-[A-Z|0-9]{16}$")

type service struct {
	service settings.CRUDService[*appdetection.Settings]
	client  rest.Client
}

func (me *service) List() (api.Stubs, error) {
	return me.service.List()
}

func (me *service) Get(id string, v *appdetection.Settings) error {
	return me.service.Get(id, v)
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(v *appdetection.Settings) (*api.Stub, error) {
	var stub *api.Stub
	var err error
	checkedAppExists := false
	maxIterations := 50
	for iteration := 0; iteration < maxIterations; iteration++ {
		stub, err = me.service.Create(v)
		if err == nil {
			break
		}
		matches := NotAValidDataSourceRegexp.FindStringSubmatch(err.Error())
		if len(matches) == 0 {
			break
		}
		applicationID := matches[1]
		if !ApplicationIDRegexp.MatchString(applicationID) {
			break
		}
		if !checkedAppExists {
			cerr := me.client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s", url.PathEscape(applicationID))).Expect(200).Finish()
			if cerr != nil {
				break
			}
			checkedAppExists = true
		}
		time.Sleep(time.Duration(500*iteration) * time.Millisecond)
	}
	return stub, err
}

func (me *service) Update(id string, v *appdetection.Settings) error {
	return me.service.Update(id, v)
}

func (me *service) Delete(id string) error {
	return me.service.Delete(id)
}

func Duplicates(service settings.RService[*appdetection.Settings], v *appdetection.Settings) (*api.Stub, error) {
	if settings.RejectDuplicate("dynatrace_application_detection_rule_v2") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*appdetection.Settings)
			if v.Matcher == config.Matcher && v.Pattern == config.Pattern {
				return nil, fmt.Errorf("Application detection rule with filter already exists")
			}
		}
	} else if settings.HijackDuplicate("dynatrace_application_detection_rule_v2") {
		var err error
		var stubs api.Stubs
		if stubs, err = service.List(); err != nil {
			return nil, err
		}
		for _, stub := range stubs {
			config := stub.Value.(*appdetection.Settings)
			if v.Matcher == config.Matcher && v.Pattern == config.Pattern {
				return stub, nil
			}
		}
	}
	return nil, nil
}
