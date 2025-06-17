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

package customtags

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customtags/list"
	customtags "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/customtags/settings"
)

var ErrZeroMatched = os.Getenv("DYNATRACE_TAGS_ERR_ZERO_MATCHED") == "true"

func Service(credentials *rest.Credentials) settings.CRUDService[*customtags.Settings] {
	return &service{credentials: credentials}
}

type service struct {
	credentials *rest.Credentials
}

var entityIdSelectorRegexp = regexp.MustCompile("entityId\\((.*)\\)")

func (me *service) Get(ctx context.Context, selector string, v *customtags.Settings) (err error) {
	client := rest.HybridClient(me.credentials)
	if err = client.Get(ctx, fmt.Sprintf("/api/v2/tags?entitySelector=%s&from=now-3y&to=now", url.QueryEscape(selector)), 200).Finish(v); err != nil {
		return err
	}
	v.EntitySelector = selector
	if entityIdSelectorRegexp.MatchString(selector) {
		var response = struct {
			Entities []struct {
				DisplayName string `json:"displayName"`
			}
		}{}
		if err = client.Get(ctx, fmt.Sprintf("/api/v2/entities?entitySelector=%s&from=now-3y&to=now", url.QueryEscape(selector)), 200).Finish(&response); err == nil {
			if len(response.Entities) > 0 {
				v.ExportName = response.Entities[0].DisplayName
			}
		}
	}

	return nil
}

func (me *service) SchemaID() string {
	return "v2:environment:customtags"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return list.List(ctx, rest.HybridClient(me.credentials))
}

func (me *service) Validate(ctx context.Context, v *customtags.Settings) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *customtags.Settings) (*api.Stub, error) {
	if err := me.Update(ctx, v.EntitySelector, v); err != nil {
		return nil, err
	}
	return &api.Stub{ID: v.EntitySelector, Name: v.EntitySelector}, nil
}

func (me *service) Update(ctx context.Context, id string, v *customtags.Settings) error {
	var err error

	var settingsObj customtags.Settings
	client := rest.HybridClient(me.credentials)
	if err = client.Post(ctx, fmt.Sprintf("/api/v2/tags?entitySelector=%s&from=now-3y&to=now", url.QueryEscape(v.EntitySelector)), v, 200).Finish(&settingsObj); err != nil {
		return err
	}
	if ErrZeroMatched && settingsObj.MatchedEntities == 0 {
		return rest.Error{Message: fmt.Sprintf("No entities matching the selector '%s' were found within the past three years.", v.EntitySelector)}
	}

	v.MatchedEntities = settingsObj.MatchedEntities

	return nil
}

func (me *service) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (me *service) DeleteValue(ctx context.Context, v *customtags.Settings) error {
	client := rest.HybridClient(me.credentials)
	for _, tag := range v.Tags {
		if tag.Value == nil || len(*tag.Value) == 0 {
			if err := client.Delete(ctx, fmt.Sprintf("/api/v2/tags?key=%s&entitySelector=%s", url.QueryEscape(tag.Key), url.QueryEscape(v.EntitySelector)), 200).Finish(); err != nil {
				return err
			}
		} else {
			if err := client.Delete(ctx, fmt.Sprintf("/api/v2/tags?key=%s&value=%s&entitySelector=%s", url.QueryEscape(tag.Key), url.QueryEscape(*tag.Value), url.QueryEscape(v.EntitySelector)), 200).Finish(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *service) New() *customtags.Settings {
	return new(customtags.Settings)
}
