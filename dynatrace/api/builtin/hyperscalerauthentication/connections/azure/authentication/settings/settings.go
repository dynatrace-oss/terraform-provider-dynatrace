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

package connectionauthentication

import (
	"context"
	"errors"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/hyperscalerauthentication/connections/azure"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Settings struct {
	Name              string
	AzureConnectionID string
	ApplicationID     *string // Application (client) ID of your app registered in Microsoft Azure App registrations
	DirectoryID       *string // Directory (tenant) ID of Microsoft Entra ID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"azure_connection_id": {
			Type:             schema.TypeString,
			Description:      "The ID of a `dynatrace_azure_connection` resource instance for which to define the Azure Authentication",
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"application_id": {
			Type:             schema.TypeString,
			Description:      "Application (client) ID of your app registered in Microsoft Azure App registrations",
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"directory_id": {
			Type:             schema.TypeString,
			Description:      "Directory (tenant) ID of Microsoft Entra ID",
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
	}
}

func (me *Settings) Timeouts() *schema.ResourceTimeout {
	return &schema.ResourceTimeout{
		Create: schema.DefaultTimeout(azure.DefaultTimeout),
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"azure_connection_id": me.AzureConnectionID,
		"application_id":      me.ApplicationID,
		"directory_id":        me.DirectoryID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"azure_connection_id": &me.AzureConnectionID,
		"application_id":      &me.ApplicationID,
		"directory_id":        &me.DirectoryID,
	})
}

func (me *Settings) CustomizeDiff(ctx context.Context, rd *schema.ResourceDiff, i any) error {
	// try to check if an attempt is being made to change application and directory IDs for a connection where they were already set.
	// In this case, we want to block the change by returning an error.

	if rd.HasChange("azure_connection_id") {
		return nil
	}

	err := errors.Join(
		checkIfKeyHasAlreadyBeenSet(rd, "application_id"),
		checkIfKeyHasAlreadyBeenSet(rd, "directory_id"),
	)
	if err != nil {
		return fmt.Errorf("%w. To fix this, destroy and recreate the underlying `dynatrace_azure_connection` resource", err)
	}

	return nil
}

func checkIfKeyHasAlreadyBeenSet(rd *schema.ResourceDiff, key string) error {
	if !rd.HasChange(key) {
		return nil
	}
	oldVal, _ := rd.GetChange(key)
	if oldVal != "" {
		return fmt.Errorf("the '%s' property is has already been set and cannot be changed", key)
	}
	return nil
}
