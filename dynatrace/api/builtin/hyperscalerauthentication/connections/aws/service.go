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
	"context"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	awsconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/aws/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaID = "builtin:hyperscaler-authentication.connections.aws"
const SchemaVersion = "0.0.15"

func Service(credentials *rest.Credentials) settings.CRUDService[*awsconnection.Settings] {
	return &service{settings20.Service[*awsconnection.Settings](credentials, SchemaID, SchemaVersion)}
}

type service struct {
	service settings.CRUDService[*awsconnection.Settings]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *awsconnection.Settings) error {
	if err := me.service.Get(ctx, id, v); err != nil {
		return err
	}
	return nil
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

func (me *service) Create(ctx context.Context, v *awsconnection.Settings) (*api.Stub, error) {
	return me.service.Create(ctx, v)
}

func (me *service) Update(ctx context.Context, id string, v *awsconnection.Settings) error {
	return nil

	///
	/// An update "should" never happen. All attributes and subresources
	/// are flagged as "forceNew"
	///
	// remoteValue := awsconnection.Settings{}
	// if err := me.Get(ctx, id, &remoteValue); err != nil {
	// 	return err
	// }
	// if remoteValue.AWSRoleBasedAuthentication != nil && v.AWSRoleBasedAuthentication != nil {
	// 	v.AWSRoleBasedAuthentication.RoleARN = remoteValue.AWSRoleBasedAuthentication.RoleARN
	// } else if remoteValue.AWSWebIdentity != nil && v.AWSWebIdentity != nil {
	// 	v.AWSWebIdentity.RoleARN = remoteValue.AWSWebIdentity.RoleARN
	// } else {
	// 	return nil
	// }
	// return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}
