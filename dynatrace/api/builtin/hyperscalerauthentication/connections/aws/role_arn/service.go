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

package role_arn

import (
	"context"
	"errors"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	awsconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws"
	role_arn "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws/role_arn/settings"
	awsconnection_settings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaID = "builtin:hyperscaler-authentication.connections.aws"
const SchemaVersion = "0.0.15"

func Service(credentials *rest.Credentials) settings.CRUDService[*role_arn.Settings] {
	return &service{
		service:     settings20.Service[*role_arn.Settings](credentials, SchemaID, SchemaVersion),
		connService: awsconnection.Service(credentials),
	}
}

type service struct {
	service     settings.CRUDService[*role_arn.Settings]
	connService settings.CRUDService[*awsconnection_settings.Settings]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.connService.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *role_arn.Settings) error {
	connValue := awsconnection_settings.Settings{}
	if err := me.connService.Get(ctx, id, &connValue); err != nil {
		return err
	}
	if connValue.AWSRoleBasedAuthentication != nil {
		v.RoleARN = connValue.AWSRoleBasedAuthentication.RoleARN
	} else if connValue.AWSWebIdentity != nil {
		v.RoleARN = connValue.AWSWebIdentity.RoleARN
	}
	v.AWSConnectionID = id
	v.Name = connValue.Name
	return nil
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(ctx context.Context, v *role_arn.Settings) (*api.Stub, error) {
	connValue := awsconnection_settings.Settings{}
	if err := me.connService.Get(ctx, v.AWSConnectionID, &connValue); err != nil {
		return nil, err
	}
	if connValue.AWSRoleBasedAuthentication != nil {
		connValue.AWSRoleBasedAuthentication.RoleARN = v.RoleARN
	} else if connValue.AWSWebIdentity != nil {
		connValue.AWSWebIdentity.RoleARN = v.RoleARN
	}
	if err := me.connService.Update(ctx, v.AWSConnectionID, &connValue); err != nil {
		return nil, err
	}

	return &api.Stub{ID: v.AWSConnectionID, Name: v.AWSConnectionID}, nil
}

func (me *service) Update(_ context.Context, _ string, _ *role_arn.Settings) error {
	return errors.New("update not supported: This resource is immutable after creation. Changes require replacement")
}

func (me *service) Delete(_ context.Context, _ string) error {
	return nil
	// Doesn't work right now - even updating to an empty roleARN errors out
	// return me.Update(ctx, id, &role_arn.Settings{AWSConnectionID: id, RoleARN: ""})
}
