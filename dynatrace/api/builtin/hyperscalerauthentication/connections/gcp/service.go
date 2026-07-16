/**
* @license
* Copyright 2026 Dynatrace LLC
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

package gcp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	serviceSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp/settings"
	retrycommon "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/retry"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
)

const SchemaVersion = "0.0.5"
const SchemaID = "builtin:hyperscaler-authentication.connections.gcp"

// RetryableAuthenticationErrorMessage is the substring of the constraint-violation message the
// Dynatrace API returns while the service account impersonation cannot yet authenticate (e.g. the
// IAM policy change on the GCP side has not propagated yet). Create retries for as long as the API
// keeps returning this message. The acceptance test asserts on the same constant, binding the
// retried message and the asserted message together: if the API ever changes the wording, the test
// fails and forces this constant to be updated in lockstep.
const RetryableAuthenticationErrorMessage = "GCP authentication failed"

func Service(clientSet rest.ClientSet) (settings.CRUDService[*serviceSettings.Settings], error) {
	svc, err := settings20.Service[*serviceSettings.Settings](clientSet, SchemaID, SchemaVersion)
	if err != nil {
		return nil, err
	}
	return &service{
		service: svc,
	}, nil
}

type service struct {
	service settings.CRUDService[*serviceSettings.Settings]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.service.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *serviceSettings.Settings) error {
	return me.service.Get(ctx, id, v)
}

func (me *service) Validate(ctx context.Context, v *serviceSettings.Settings) error {
	if validator, ok := me.service.(settings.Validator[*serviceSettings.Settings]); ok {
		return validator.Validate(ctx, v)
	}

	return settings.ErrValidatorNotImplemented(me.SchemaID())
}

func (me *service) SchemaID() string {
	return me.service.SchemaID()
}

// Create is wrapped in a retry: the service account impersonation relies on the provided
// service account ID, and the corresponding IAM policy changes on the GCP side may not have
// propagated yet. We retry until they eventually do (or the create timeout elapses).
func (me *service) Create(ctx context.Context, v *serviceSettings.Settings) (*api.Stub, error) {
	var stub *api.Stub
	err := retry.RetryContext(ctx, retrycommon.DurationUntilDeadlineOrDefault(ctx, serviceSettings.DefaultCreateTimeout), func() *retry.RetryError {
		var err error
		stub, err = me.service.Create(ctx, v)

		return settings.ClassifyConstraintRetryError(err, RetryableAuthenticationErrorMessage)

	})
	if err != nil {
		return nil, err
	}
	return stub, nil
}

// Update only ever changes the connection name; service_account_id and the rest of the
// impersonation config are ForceNew, so they go through Create/Delete instead. A rename does
// not depend on IAM propagation, so no retry is needed here.
func (me *service) Update(ctx context.Context, id string, v *serviceSettings.Settings) error {
	return me.service.Update(ctx, id, v)
}

func (me *service) Delete(ctx context.Context, id string) error {
	return me.service.Delete(ctx, id)
}
