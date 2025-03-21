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

package directshares

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	directshares "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/documents/directshares/settings"
	httplog "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest/logging"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/auth"
	directshare "github.com/dynatrace-oss/terraform-provider-dynatrace/monaco/pkg/client/document/directshares"
)

func Service(credentials *rest.Credentials) settings.CRUDService[*directshares.DirectShare] {
	return &service{credentials}
}

type service struct {
	credentials *rest.Credentials
}

func (me *service) client(ctx context.Context) *directshare.Client {
	httplog.InstallRoundTripper()
	httpClient := auth.NewOAuthClient(ctx, auth.OauthCredentials{
		ClientID:     me.credentials.OAuth.ClientID,
		ClientSecret: me.credentials.OAuth.ClientSecret,
		TokenURL:     me.credentials.OAuth.TokenURL,
	})
	return directshare.NewClient(me.credentials.OAuth.EnvironmentURL, httpClient)
}

func (me *service) Get(ctx context.Context, id string, v *directshares.DirectShare) (err error) {
	var result *directshare.Response
	if result, err = me.client(ctx).GET(ctx, directshare.DirectShares, id); err != nil {
		return err
	}

	v.DocumentId = result.DocumentId
	v.ID = result.ID
	v.Access = strings.Join(result.Access, "-")
	v.Recipients = result.Recipients

	return nil
}

func (me *service) SchemaID() string {
	return "document:direct-shares"
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	var stubs api.Stubs
	return stubs, nil // not implemented

}

func (me *service) Validate(v *directshares.DirectShare) error {
	return nil // no endpoint for that
}

func (me *service) Create(ctx context.Context, v *directshares.DirectShare) (stub *api.Stub, err error) {
	var id string

	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	if id, err = me.client(ctx).INSERT(ctx, directshare.DirectShares, data); err != nil {
		return nil, err
	}
	return &api.Stub{ID: id}, nil
}

func (me *service) Update(ctx context.Context, id string, v *directshares.DirectShare) (err error) {

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return me.client(ctx).UPDATE(ctx, directshare.DirectShares, id, data)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.client(ctx).DELETE(ctx, directshare.DirectShares, id)
}

func (me *service) New() *directshares.DirectShare {
	return new(directshares.DirectShare)
}
