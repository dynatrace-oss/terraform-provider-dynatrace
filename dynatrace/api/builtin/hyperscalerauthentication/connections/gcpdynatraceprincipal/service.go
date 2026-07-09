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

package gcpdynatraceprincipal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp"
	gcpconnection "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp/settings"
	gcpprincipalsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcpdynatraceprincipal/settings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
)

const SchemaVersion = "0.0.3"
const SchemaID = "builtin:hyperscaler-authentication.connections.gcp-dynatrace-principal"
const principalCreationTimeout = 5 * time.Minute

var errPrincipalNotProvisionedYet = errors.New("no Dynatrace GCP Principal found yet")

var ValidConnection = gcpconnection.Settings{
	Name: "temp",
	Type: "serviceAccountImpersonation",
	ServiceAccountImpersonation: new(gcpconnection.ServiceAccountImpersonation{
		Consumers: []gcpconnection.ConsumersOfServiceAccountImpersonation{"SVC:com.dynatrace.bo"},
	}),
}

func Service(credentials *config.ProviderConfiguration) settings.CRUDService[*gcpprincipalsettings.Settings] {
	return &service{
		principalService:  settings20.Service[*gcpprincipalsettings.Settings](credentials, SchemaID, SchemaVersion),
		connectionService: gcp.Service(credentials),
	}
}

type service struct {
	principalService  settings.CRUDService[*gcpprincipalsettings.Settings]
	connectionService settings.CRUDService[*gcpconnection.Settings]
}

func (me *service) List(ctx context.Context) (api.Stubs, error) {
	return me.principalService.List(ctx)
}

func (me *service) Get(ctx context.Context, id string, v *gcpprincipalsettings.Settings) error {
	return me.principalService.Get(ctx, id, v)
}

func (me *service) SchemaID() string {
	return me.principalService.SchemaID()
}

func (me *service) ensurePrincipalExists(ctx context.Context) (*api.Stub, error) {
	stubs, err := me.principalService.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(stubs) == 0 {
		// Provisioning is triggered by a validateOnly call against the GCP connection schema, not
		// the principal schema. As a consequence the caller needs write access on the GCP connection
		// schema (builtin:hyperscaler-authentication.connections.gcp) to provision the principal.
		// Write access on the principal schema is not required at all — a caller denied write access
		// on the principal schema can still trigger its creation here. This is a Dynatrace API
		// characteristic; we cannot work around it here.
		if v, ok := me.connectionService.(settings.Validator[*gcpconnection.Settings]); ok {
			if err = v.Validate(ctx, &ValidConnection); err != nil {
				return nil, fmt.Errorf("failed to trigger Dynatrace GCP Principal creation: %w", err)
			}
		} else {
			return nil, settings.ErrValidatorNotImplemented(gcp.SchemaID)
		}

		err = retry.RetryContext(ctx, principalCreationTimeout, func() *retry.RetryError {
			var listErr error
			stubs, listErr = me.principalService.List(ctx)
			if listErr != nil {
				return retry.NonRetryableError(fmt.Errorf("Dynatrace GCP Principal provisioning failed: %w", listErr))
			}
			if len(stubs) == 0 {
				return retry.RetryableError(errPrincipalNotProvisionedYet)
			}
			return nil
		})
		if err != nil {
			if errors.Is(err, errPrincipalNotProvisionedYet) {
				return nil, fmt.Errorf("Dynatrace GCP Principal was not provisioned within %v", principalCreationTimeout)
			}
			return nil, err
		}
	}

	// There can only be one principal, so we take the first (and only one) in the list.
	return stubs[0], nil
}

// Create ignores its settings argument: the Dynatrace GCP Principal takes no user input. It is a
// singleton that is provisioned by Dynatrace as a side effect of submitting ValidConnection to the
// GCP connection's validate endpoint. If the principal already exists this is a no-op that simply
// returns the existing stub; otherwise creation is triggered and ensurePrincipalExists polls until
// the principal appears (or principalCreationTimeout elapses).
func (me *service) Create(ctx context.Context, _ *gcpprincipalsettings.Settings) (*api.Stub, error) {
	return me.ensurePrincipalExists(ctx)
}

func (me *service) Update(_ context.Context, _ string, _ *gcpprincipalsettings.Settings) error {
	// Noop
	return nil
}

func (me *service) Delete(_ context.Context, _ string) error {
	// Noop
	return nil
}
