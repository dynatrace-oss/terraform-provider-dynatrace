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
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/cache"

	webservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web"
	dataprivacy "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/dataprivacy/settings"
	web "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/applications/web/settings"
)

const SchemaID = "v1:config:applications:web:data-privacy"

func Service(credentials *rest.Credentials) settings.CRUDService[*dataprivacy.ApplicationDataPrivacy] {
	return &service{
		schemaID:      SchemaID,
		client:        rest.APITokenClient(credentials),
		webAppService: cache.CRUD(webservice.Service(credentials), true)}
}

type service struct {
	schemaID      string
	client        rest.Client
	webAppService settings.CRUDService[*web.Application]
}

func extractApplicationID(id string) string {
	if strings.HasSuffix(id, "-data-privacy") {
		return strings.TrimSuffix(id, "-data-privacy")
	} else if strings.HasPrefix(id, "DATA-PRIVACY-") {
		return strings.TrimPrefix(id, "DATA-PRIVACY-")
	}
	return id
}

func (me *service) Get(ctx context.Context, id string, v *dataprivacy.ApplicationDataPrivacy) error {
	id = extractApplicationID(id)

	req := me.client.Get(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/dataPrivacy", url.PathEscape(id)), 200)

	if err := req.Finish(v); err != nil {
		return err
	}

	stubs, err := me.webAppService.List(ctx)
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

func (me *service) Update(ctx context.Context, id string, v *dataprivacy.ApplicationDataPrivacy) error {
	id = extractApplicationID(id)
	err := me.client.Put(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/dataPrivacy", id), v, 201, 204).Finish()
	if err != nil && strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return nil
	}
	return err
}

func (me *service) Delete(ctx context.Context, id string) error {
	id = extractApplicationID(id)
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

	return me.Update(ctx, id, &settings)
}

func (me *service) Validate(ctx context.Context, v *dataprivacy.ApplicationDataPrivacy) error {
	id := *v.WebApplicationID
	id = extractApplicationID(id)
	err := me.client.Post(ctx, fmt.Sprintf("/api/config/v1/applications/web/%s/dataPrivacy/validator", id), v, 204).Finish()
	if err != nil && strings.HasPrefix(err.Error(), "No Content (PUT)") {
		return nil
	}
	return err
}

func (me *service) Create(ctx context.Context, v *dataprivacy.ApplicationDataPrivacy) (*api.Stub, error) {
	if err := me.Update(ctx, *v.WebApplicationID, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: *v.WebApplicationID + "-data-privacy"}, nil
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var err error
	var stubs api.Stubs

	if stubs, err = me.webAppService.List(ctx); err != nil {
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
