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

package gcpdynatraceprincipal

import (
	"context"
	"errors"
	"fmt"
	"time"

	gcpdynatraceprincipal "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/hyperscalerauthentication/connections/gcpdynatraceprincipal/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp"
	gcpsettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/gcp/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/services/settings20"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	SchemaVersion = "0.0.3"
	SchemaID      = "builtin:hyperscaler-authentication.connections.gcp-dynatrace-principal"

	principalCreationTimeout = 5 * time.Minute
)

var errPrincipalNotProvisionedYet = errors.New("no Dynatrace GCP Principal found yet")

var ValidConnection = gcpsettings.Settings{
	Name: "temp",
	Type: "serviceAccountImpersonation",
	ServiceAccountImpersonation: new(gcpsettings.ServiceAccountImpersonation{
		Consumers: []gcpsettings.ConsumersOfServiceAccountImpersonation{"SVC:com.dynatrace.bo"},
	}),
}

func ensurePrincipalExists(ctx context.Context, credentials *rest.Credentials) (*api.Stub, error) {
	listService := settings20.Service[*gcpdynatraceprincipal.Settings](credentials, SchemaID, SchemaVersion)

	stubs, err := listService.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(stubs) == 0 {
		if v, ok := gcp.Service(credentials).(settings.Validator[*gcpsettings.Settings]); ok {
			if err = v.Validate(ctx, &ValidConnection); err != nil {
				return nil, fmt.Errorf("failed to trigger Dynatrace GCP Principal creation. Please reach out to support: %w", err)
			}
		} else {
			return nil, settings.ErrValidatorNotImplemented(gcp.SchemaID)
		}

		err = retry.RetryContext(ctx, principalCreationTimeout, func() *retry.RetryError {
			var listErr error
			stubs, listErr = listService.List(ctx)
			if listErr != nil {
				return retry.NonRetryableError(fmt.Errorf("Dynatrace GCP Principal provisioning failed. Please reach out to support: %w", listErr))
			}
			if len(stubs) == 0 {
				return retry.RetryableError(errPrincipalNotProvisionedYet)
			}
			return nil
		})
		if err != nil {
			if errors.Is(err, errPrincipalNotProvisionedYet) {
				return nil, fmt.Errorf("Dynatrace GCP Principal was not provisioned within %v. Please reach out to support", principalCreationTimeout)
			}
			return nil, err
		}
	}

	// There can only be one principal, so we take the first (and only one) in the list.
	return stubs[0], nil
}

// Attribute keys
const (
	attrPrincipal = "principal"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: logging.EnableDSCtx(DataSourceRead),
		Description: "Returns the Dynatrace GCP Principal",
		Schema: map[string]*schema.Schema{
			attrPrincipal: {
				Type:        schema.TypeString,
				Description: "Dynatrace GCP Principal",
				Computed:    true,
			},
		},
	}
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	principalStub, err := ensurePrincipalExists(ctx, creds)
	if err != nil {
		return diag.FromErr(err)
	}

	principal, ok := principalStub.Value.(*gcpdynatraceprincipal.Settings)
	if !ok || principal == nil {
		return diag.FromErr(fmt.Errorf("unexpected value type returned for Dynatrace GCP Principal. Please reach out to support"))
	}

	d.SetId(principalStub.ID)
	if err := d.Set(attrPrincipal, principal.Principal); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
