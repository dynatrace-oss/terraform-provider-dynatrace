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

package dataprivacy

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	webservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web"
	dataprivacy "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/dataprivacy/settings"
	web "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings"
)

const SchemaID = "v1:config:applications:web:data-privacy"

func Service(credentials *settings.Credentials) settings.CRUDService[*dataprivacy.ApplicationDataPrivacy] {
	return &service{
		schemaID:      SchemaID,
		client:        rest.DefaultClient(credentials.URL, credentials.Token),
		webAppService: cache.CRUD(webservice.Service(credentials), true)}
}

type service struct {
	schemaID      string
	client        rest.Client
	webAppService settings.CRUDService[*web.Application]
}

func (me *service) Get(id string, v *dataprivacy.ApplicationDataPrivacy) error {
	id = strings.TrimSuffix(id, "-data-privacy")
	req := me.client.Get(fmt.Sprintf("/api/config/v1/applications/web/%s/dataPrivacy", url.PathEscape(id)), 200)

	if err := req.Finish(v); err != nil {
		return err
	}

	stubs, err := me.webAppService.List()
	if err != nil {
		return err
	}
	for _, stub := range stubs {
		if stub.ID == id {
			v.Name = "Data Privacy Settings for " + stub.Name
			v.WebApplicationID = &id
			break
		}
	}
	// application := web.Application{}
	// if err := me.webAppService.Get(id, &application); err != nil {
	// 	return err
	// }
	return nil
}

func (me *service) Update(id string, v *dataprivacy.ApplicationDataPrivacy) error {
	id = strings.TrimSuffix(id, "-data-privacy")
	err := me.client.Put(fmt.Sprintf("/api/config/v1/applications/web/%s/dataPrivacy", id), v, 201, 204).Finish()
	if err != nil && strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return nil
	}
	return err
}

func (me *service) Delete(id string) error {
	id = strings.TrimSuffix(id, "-data-privacy")
	settings := dataprivacy.ApplicationDataPrivacy{
		DataCaptureOptInEnabled:         false,
		PersistentCookieForUserTracking: false,
		DoNotTrackBehaviour:             dataprivacy.DoNotTrackBehaviours.CaptureAnonymized,
		SessionReplayDataPrivacy: &dataprivacy.SessionReplayDataPrivacySettings{
			OptIn:             false,
			URLExclusionRules: []string{},
			ContentMaskingSettings: &dataprivacy.ContentMaskingSettings{
				RecordingMaskingSettingsVersion: 2,
				RecordingMaskingSettings: &dataprivacy.SessionReplayMaskingSetting{
					Preset: dataprivacy.MaskingPresets.MaskAll,
					Rules:  dataprivacy.MaskingRules{},
				},
				PlaybackMaskingSettings: &dataprivacy.SessionReplayMaskingSetting{
					Preset: "MASK_ALL",
					Rules:  dataprivacy.MaskingRules{},
				},
			},
		},
	}

	return me.Update(id, &settings)
}

func (me *service) Validate(v *dataprivacy.ApplicationDataPrivacy) error {
	id := *v.WebApplicationID
	id = strings.TrimSuffix(id, "-data-privacy")
	err := me.client.Post(fmt.Sprintf("/api/config/v1/applications/web/%s/dataPrivacy/validator", id), v, 204).Finish()
	if err != nil && strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return nil
	}
	return err
}

func (me *service) Create(v *dataprivacy.ApplicationDataPrivacy) (*settings.Stub, error) {
	if err := me.Update(*v.WebApplicationID, v); err != nil {
		return nil, err
	}
	return &settings.Stub{ID: *v.WebApplicationID + "-data-privacy"}, nil
}

func (me *service) List() (settings.Stubs, error) {
	var err error
	var stubs settings.Stubs

	if stubs, err = me.webAppService.List(); err != nil {
		return nil, err
	}
	for _, stub := range stubs {
		stub.Name = "Data Privacy Settings for " + stub.Name
		stub.ID = stub.ID + "-data-privacy"
	}
	return stubs.ToStubs(), nil
}

func (me *service) SchemaID() string {
	return me.schemaID
}
