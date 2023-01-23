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

package externalvault

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	SourceAuthMethod                          SourceAuthMethod `json:"sourceAuthMethod"` // Defines the actual set of fields depending on the value. See one of the following objects: \n\n* `HASHICORP_VAULT_APPROLE` -> HashicorpApproleConfig \n* `HASHICORP_VAULT_CERTIFICATE` -> HashicorpCertificateConfig \n* `AZURE_KEY_VAULT_CLIENT_SECRET` -> AzureClientSecretConfig \n
	VaultURL                                  *string          `json:"vaultUrl,omitempty"`
	UsernameSecretName                        *string          `json:"usernameSecretName,omitempty"`
	PasswordSecretName                        *string          `json:"passwordSecretName,omitempty"`
	TokenSecretName                           *string          `json:"tokenSecretName,omitempty"`
	CredentialsUsedForExternalSynchronization []string         `json:"credentialsUsedForExternalSynchronization,omitempty"`

	// HashicorpApproleConfig
	PathtoCredentials *string `json:"pathToCredentials,omitempty"`
	RoleID            *string `json:"roleId,omitempty"`
	SecretID          *string `json:"secretId,omitempty"` // The ID of Credentials within the Certificate Vault holding the secret id
	VaultNameSpace    *string `json:"vaultNamespace,omitempty"`
	// HashicorpCertificateConfig
	Certificate *string `json:"certificate,omitempty"` // The ID of Credentials within the Certificate Vault holding the certificate
	// AzureClientSecret
	TenantID     *string `json:"tenantId,omitempty"`     // Tenant (directory) ID of Azure application in Azure Active Directory which has permission to access secrets in Azure Key Vault
	ClientID     *string `json:"clientId,omitempty"`     // Client (application) ID of Azure application in Azure Active Directory which has permission to access secrets in Azure Key Vault
	ClientSecret *string `json:"clientSecret,omitempty"` // Client secret generated for Azure application in Azure Active Directory used for proving identity when requesting a token used later for accessing secrets in Azure Key Vault
}

func (me *Config) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_secret": {
			Type:        schema.TypeString,
			Description: "Required for Azure Client Secret. No further documentation available",
			Optional:    true,
		},
		"clientid": {
			Type:        schema.TypeString,
			Description: "Required for Azure Client Secret. No further documentation available",
			Optional:    true,
		},
		"tenantid": {
			Type:        schema.TypeString,
			Description: "Required for Azure Client Secret. No further documentation available",
			Optional:    true,
		},
		"certificate": {
			Type:        schema.TypeString,
			Description: "Required for Hashicorp Certificate. The ID of Credentials within the Certificate Vault holding the certificate",
			Optional:    true,
		},
		"vault_namespace": {
			Type:        schema.TypeString,
			Description: "Required for Hashicorp App Role. No further documentation available",
			Optional:    true,
		},
		"secretid": {
			Type:        schema.TypeString,
			Description: "Required for Hashicorp App Role. The ID of Credentials within the Certificate Vault holding the secret id",
			Optional:    true,
		},
		"roleid": {
			Type:        schema.TypeString,
			Description: "Required for Hashicorp App Role. No further documentation available",
			Optional:    true,
		},
		"path_to_credentials": {
			Type:        schema.TypeString,
			Description: "Required for Hashicorp App Role or Hashicorp Certificate. No further documentation available",
			Optional:    true,
		},
		"vault_url": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true,
		},
		"username_secret_name": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true,
		},
		"password_secret_name": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true,
		},
		"token_secret_name": {
			Type:        schema.TypeString,
			Description: "No documentation available",
			Optional:    true,
		},
		"credentials_used_for_external_synchronization": {
			Type:        schema.TypeSet,
			Description: "No documentation available",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
		},
	}
}

func (me *Config) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("tenantid", me.TenantID); err != nil {
		return err
	}
	if err := properties.Encode("clientid", me.ClientID); err != nil {
		return err
	}
	if err := properties.Encode("client_secret", me.ClientSecret); err != nil {
		return err
	}

	if err := properties.Encode("path_to_credentials", me.PathtoCredentials); err != nil {
		return err
	}
	if err := properties.Encode("roleid", me.RoleID); err != nil {
		return err
	}
	if err := properties.Encode("secretid", me.SecretID); err != nil {
		return err
	}
	if err := properties.Encode("vault_namespace", me.VaultNameSpace); err != nil {
		return err
	}
	if err := properties.Encode("certificate", me.Certificate); err != nil {
		return err
	}

	if err := properties.Encode("vault_url", me.VaultURL); err != nil {
		return err
	}
	if err := properties.Encode("username_secret_name", me.UsernameSecretName); err != nil {
		return err
	}
	if err := properties.Encode("password_secret_name", me.PasswordSecretName); err != nil {
		return err
	}
	if err := properties.Encode("token_secret_name", me.TokenSecretName); err != nil {
		return err
	}
	if err := properties.Encode("credentials_used_for_external_synchronization", me.CredentialsUsedForExternalSynchronization); err != nil {
		return err
	}

	return nil
}

func (me *Config) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("client_secret"); ok {
		me.ClientSecret = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("clientid"); ok {
		me.ClientID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("tenantid"); ok {
		me.TenantID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("certificate"); ok {
		me.Certificate = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("vault_namespace"); ok {
		me.VaultNameSpace = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("secretid"); ok {
		me.SecretID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("roleid"); ok {
		me.RoleID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("path_to_credentials"); ok {
		me.PathtoCredentials = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("vault_url"); ok {
		me.VaultURL = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("username_secret_name"); ok {
		me.UsernameSecretName = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("password_secret_name"); ok {
		me.PasswordSecretName = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("token_secret_name"); ok {
		me.TokenSecretName = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("credentials_used_for_external_synchronization"); ok {
		me.CredentialsUsedForExternalSynchronization = []string{}
		for _, elem := range value.(*schema.Set).List() {
			me.CredentialsUsedForExternalSynchronization = append(me.CredentialsUsedForExternalSynchronization, elem.(string))
		}
	}
	if me.ClientID != nil || me.ClientSecret != nil || me.TenantID != nil {
		me.SourceAuthMethod = SourceAuthMethods.AzureKeyVaultClientSecret
	} else if me.RoleID != nil {
		me.SourceAuthMethod = SourceAuthMethods.HashicorpVaultAppRole
	} else if me.Certificate != nil {
		me.SourceAuthMethod = SourceAuthMethods.HashicoprVaultCertificate
	}
	return nil
}
