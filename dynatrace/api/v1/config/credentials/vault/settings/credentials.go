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

package vault

import (
	"encoding/json"
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/credentials/vault/settings/externalvault"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Credentials struct {
	Name                   string                `json:"name"`                             // The name of the credentials set.
	Description            *string               `json:"description,omitempty"`            // A short description of the credentials set..
	OwnerAccessOnly        bool                  `json:"ownerAccessOnly"`                  // The credentials set is available to every user (`false`) or to owner only (`true`).
	Scope                  Scope                 `json:"scope"`                            // The scope of the credentials set
	Type                   CredentialsType       `json:"type"`                             // Defines the actual set of fields depending on the value. See one of the following objects: \n\n* `CERTIFICATE` -> CertificateCredentials \n* `PUBLIC_CERTIFICATE` -> PublicCertificateCredentials \n* `USERNAME_PASSWORD` -> UserPasswordCredentials \n* `TOKEN` -> TokenCredentials \n
	Token                  *string               `json:"token,omitempty"`                  // Token in the string format.
	Password               *string               `json:"password,omitempty"`               // The password of the credential (Base64 encoded).
	Username               *string               `json:"user,omitempty"`                   // The username of the credentials set.
	Certificate            *string               `json:"certificate,omitempty"`            // The certificate in the string (Base64) format.
	CertificateFormat      *CertificateFormat    `json:"certificateFormat,omitempty"`      // The certificate format.
	ExternalVault          *externalvault.Config `json:"externalVault,omitempty"`          // Configuration for external vault synchronization
	CredentialUsageSummary UsageSummary          `json:"credentialUsageSummary,omitempty"` //The list contains summary data related to the use of credentials
}

func (me *Credentials) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the credentials set",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "A short description of the credentials set",
			Optional:    true,
		},
		"owner_access_only": {
			Type:        schema.TypeBool,
			Description: "The credentials set is available to every user (`false`) or to owner only (`true`)",
			Optional:    true,
		},
		"public": {
			Type:          schema.TypeBool,
			Description:   "For certificate authentication specifies whether it's public certificate auth (`true`) or not (`false`).",
			ConflictsWith: []string{"username", "token"},
			Optional:      true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of the credentials set. Possible values are `ALL`, `EXTENSION` and `SYNTHETIC`",
			Required:    true,
		},
		"token": {
			Type:          schema.TypeString,
			Description:   "Token in the string format. Specifying a token implies `Token Authentication`.",
			ConflictsWith: []string{"username", "password", "certificate", "format", "public"},
			Sensitive:     true,
			Optional:      true,
		},
		"username": {
			Type:          schema.TypeString,
			Description:   "The username of the credentials set.",
			ConflictsWith: []string{"token", "public", "certificate"},
			RequiredWith:  []string{"password"},
			Sensitive:     true,
			Optional:      true,
		},
		"password": {
			Type:          schema.TypeString,
			Description:   "The password of the credential.",
			ConflictsWith: []string{"token"},
			Sensitive:     true,
			Optional:      true,
		},
		"certificate": {
			Type:          schema.TypeString,
			Description:   "The certificate in the string format.",
			ConflictsWith: []string{"token", "username"},
			RequiredWith:  []string{"format", "password"},
			Optional:      true,
		},
		"format": {
			Type:          schema.TypeString,
			Description:   "The certificate format. Possible values are `PEM`, `PKCS12` and `UNKNOWN`.",
			ConflictsWith: []string{"token", "username"},
			RequiredWith:  []string{"certificate"},
			Optional:      true,
		},
		"external": {
			Type:        schema.TypeList,
			Description: "External Vault Configuration",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(externalvault.Config).Schema()},
		},
		"credential_usage_summary": {
			Type:        schema.TypeList,
			Description: "The list contains summary data related to the use of credentials",
			Optional:    true,
			MaxItems:    2,
			Elem:        &schema.Resource{Schema: new(CredentialUsageObj).Schema()},
			Deprecated:  "`credential_usage_summary` will be removed in future versions. It's not getting filled anymore, because it's runtime data",
		},
	}
}

func (me *Credentials) EnsurePredictableOrder() {
	conds := UsageSummary{}
	condStrings := sort.StringSlice{}
	for _, entry := range me.CredentialUsageSummary {
		condBytes, _ := json.Marshal(entry)
		condStrings = append(condStrings, string(condBytes))
	}
	condStrings.Sort()
	for _, condString := range condStrings {
		cond := CredentialUsageObj{}
		json.Unmarshal([]byte(condString), &cond)
		conds = append(conds, &cond)
	}
	me.CredentialUsageSummary = conds
}

func (me *Credentials) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if me.Description != nil && len(*me.Description) > 0 {
		if err := properties.Encode("description", me.Description); err != nil {
			return err
		}
	}
	if me.OwnerAccessOnly {
		if err := properties.Encode("owner_access_only", me.OwnerAccessOnly); err != nil {
			return err
		}
	}
	if err := properties.Encode("scope", string(me.Scope)); err != nil {
		return err
	}
	if me.Token != nil && len(*me.Token) > 0 {
		if err := properties.Encode("token", me.Token); err != nil {
			return err
		}
	}
	if me.Password != nil && len(*me.Password) > 0 {
		if err := properties.Encode("password", me.Password); err != nil {
			return err
		}
	}
	if me.Username != nil && len(*me.Username) > 0 {
		if err := properties.Encode("username", me.Username); err != nil {
			return err
		}
	}
	if me.Certificate != nil && len(*me.Certificate) > 0 {
		if err := properties.Encode("certificate", me.Certificate); err != nil {
			return err
		}
	}
	if me.CertificateFormat != nil && len(*me.CertificateFormat) > 0 {
		if err := properties.Encode("format", me.CertificateFormat); err != nil {
			return err
		}
	}
	if me.ExternalVault != nil {
		marshalled := hcl.Properties{}
		err := me.ExternalVault.MarshalHCL(marshalled)
		if err != nil {
			return err
		}
		properties["external"] = []any{marshalled}
	}
	// if err := properties.Encode("credential_usage_summary", me.CredentialUsageSummary); err != nil {
	// 	return err
	// }
	properties["credential_usage_summary"] = nil

	switch me.Type {
	case CredentialsTypes.PublicCertificate:
		if err := properties.Encode("public", true); err != nil {
			return err
		}
	}
	return nil
}

func (me *Credentials) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("external.#"); ok && value.(int) == 1 {
		me.ExternalVault = new(externalvault.Config)
		me.ExternalVault.UnmarshalHCL(hcl.NewDecoder(decoder, "external.0"))
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("description"); ok {
		me.Description = opt.NewString(value.(string))
		if len(*me.Description) == 0 {
			me.Description = nil
		}
	}
	if value, ok := decoder.GetOk("owner_access_only"); ok {
		me.OwnerAccessOnly = value.(bool)
	}
	if value, ok := decoder.GetOk("scope"); ok {
		me.Scope = Scope(value.(string))
	}
	if value, ok := decoder.GetOk("token"); ok {
		me.Token = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("password"); ok {
		me.Password = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("username"); ok {
		me.Username = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("certificate"); ok {
		me.Certificate = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("format"); ok {
		me.CertificateFormat = CertificateFormat(value.(string)).Ref()
	}
	if me.Username != nil {
		me.Type = CredentialsTypes.UsernamePassword
	} else if me.Token != nil {
		me.Type = CredentialsTypes.Token
	} else if me.Certificate != nil || me.CertificateFormat != nil {
		if value, ok := decoder.GetOk("public"); ok && value.(bool) {
			me.Type = CredentialsTypes.PublicCertificate
		} else {
			me.Type = CredentialsTypes.Certificate
		}
	}
	// if result, ok := decoder.GetOk("credential_usage_summary.#"); ok {
	// 	me.CredentialUsageSummary = []*CredentialUsageObj{}
	// 	for idx := 0; idx < result.(int); idx++ {
	// 		entry := new(CredentialUsageObj)
	// 		if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "credential_usage_summary", idx)); err != nil {
	// 			return err
	// 		}
	// 		me.CredentialUsageSummary = append(me.CredentialUsageSummary, entry)
	// 	}
	// }
	return nil
}

const credsNotProvided = "REST API didn't provide credential data"

func (me *Credentials) FillDemoValues() []string {
	switch me.Type {
	case CredentialsTypes.Certificate, CredentialsTypes.PublicCertificate:
		me.Certificate = opt.NewString("MIIKUQIBAzCCChcGCSqGSIb3DQEHAaCCCggEggoEMIIKADCCBLcGCSqGSIb3DQEHBqCCBKgwggSkAgEAMIIEnQYJKoZIhvcNAQcBMBwGCiqGSIb3DQEMAQYwDgQIymH8FWQ3IfACAggAgIIEcKpc+/EZAkI2MZOFZ05x5HvcVi60rtmsaxJ4WxZE1TVioKyXumqa0Vm3Z34TDNlknSZqkWDTxZghHPiJPflbfT+GG1ZqQ9oCfo7XLm5Q6/OTndJzWhrC3eIVGntVBFYe+VtBsQI2uj3wwKlgGAUiA1aVWSJfOjdBmrVCA2qfTn6rsook3tldBo87wpz/hFXftLXKnG64o1y1bleVGrCk+gsnytdIPqUKB/XLhz1+gA2HukSluIjsoGl+lelEY3221S9n1aFR+JDvMlrdt4yGvRMrKD4tpu+Em/Saah/UvkGqiNwvsCNIJZVJalmibK7KhpYbefH7Tki6SP8Qlw+uITEy4Nxcnx3PfxdEK64N+f++qYvL1tn4da9Ag5nPRgrKwp620zIH8xtmmThKbKsWlTnDvzMvwgXvRtjTD6CiTNCl11DqMKFsu8obSG0bh+Y/7iR9LrbonNz3FtUlr58OjPlpAB/qaDL4569FWUx7twe2wZxincGjz5M4m5TJCsTc4HJYZgCMbkJIBSnjNsFF7dH3NLu1QgCH3d6I/AnWEOHVHhRHjW0ThLjVKQSMBgxvgAH0Ywqfitq2sObnoSFHJAJzv6G/ue2XY03gswF2Tc31+dSKZ8jvDL689gt5mHg68tDKkna67ShPAnXhbVyYVl6pQxzBJpr791/i5AdrERaY+lohaQWZcN03ntuqUvGNbckKO/5M5AbkTRRLOdh+c3WJEA/lDChJW/0uhU85y0a92g7bvKVVGgGnbAHsZCfAd3BprC4Ub1V9fvOBtqortwylLJQv61Kw9PxzHtVmwGIS+FQJhuHi2CeOO9aSSfxgvEZcenfCiYP1PbljI1BclD2L4tl13z3IGF2TxjR+DWL+mXj8lJCS+4VlauUG93BSd93Fxr9ogyN/9iYxLrFVdEenplQSMYjV1kxgkU5sElxGYjvjkdV8zncxvhQxr5ZwdWFUOt5QR/zjJyq2qNdRtiYnm4kyet7Ednp7XESjg0D/SYcwsN2nLXOHlAvaB/8xarOoVx5tGh5LUL0uqrVPuR8yR1jrgdKAPGUUxd+xClSnWWBF6IK2QwdZglnJzPUpPeib7nvvMHy/RTCARW01dU9m5LmjqUSlhC1KBXHvtowfSvjOFwYuVWNewf3AmbQ3y0CM4bQc91gOKP9rAMeP1awFMy1p6CdqBPmPowua4nprmZpb/2IUoyFNxCTS2+b5Vl0mH2CiSjmntD3J05vboCT7rH6CdiGruR0/5RD8yA3KITS+R6HZl2P1L7JvaTOtgCGj5niIiMjSIgJj5RyI5UIdRwIg+yzECu8t5iFGOwoM0apB9oVsXRMfNdUFSgTJ/Mk1/Gpn9kLIMc0eLPc5NyAWbkIRZAwX6omRuJw6YC1LR8iZe4I8y73tyIgKOeUrl+8BxrkYkBDS70WrLsuHP8aT0pdaJ8cMFyO7GRRmEePrF9lT0liLEbGjZv/ULPlNkTTlXdQETrhzPf3tdrt+5b0bfQtc93s40iE8FYZWABepMIIFQQYJKoZIhvcNAQcBoIIFMgSCBS4wggUqMIIFJgYLKoZIhvcNAQwKAQKgggTuMIIE6jAcBgoqhkiG9w0BDAEDMA4ECPCjMDeRKs0SAgIIAASCBMjh5pysncWWvg2MORvnTIb50uaaOEl/oJhTfoXJAEVZGeiP/Sv1+YxP0wFFcrKwS2jnry8Xbw0vsumec6yo+QGshGLhJqFSrikf1oZi2F/zTPM7iBf4VUYY5AgiybHVnUU42Uh0g9mFKS6VQnaPSmeil5EtOBRFtYg5UC+1tDBw0sc/ue15uoA5UihjJm60dtGhkbxH+3T/QkgT1B0BlnHnamlpNiw0eQfKeO00m3FA23s8HgVkVvgOq32G9mB64MjexJj6b+qjhoDvNXBdRszwnDySkxbLlPBEMF5xD4OSVw57OJAr8lsTY+Ma47vIjO6zlAXQIi9vU/kfurbLATbIcOiYgDvFYuPeYZ1fo5E7Sff03oYTFOKjC/xa+oTcjA2L36vl+yKFluRYbx00NIB7BCvR7jfzX+ojpiupLODE0Yne4SJKXdaDWm1buDBWHEKCklWsYQAquPQagC/JOLSSThChcpS2xz0dvuxNfzRWy4f1NkQyD823ijTegeZkAeMBApfpAYe2yb2JMfkE6fZUmENjmY9pjXLfGEWAUQciFXQL42orYVLnU13ai6j3CEVsP30+9ZiOkaH72BDX+QnQ1h5oTL2PRT9CT8KXMrcDQPRF/blaD1Q8IG3bdPUO8X+ij7HRxsR+3llf1mSg7HAVo2nPiq1GwwNlkDOZ/aVu5P8zZi0fJrdOoWL7EjlWFcHthKCGH1829Q4fSDjkw0R/itTERhYWHxhlU8u1RXhbClzatq63UKYOBGxccVc1L5UaHVdDIaXXOhT5kYBEAtavee7c+J/UpK94fVQ3BbMjucnxT5fJqVVwqFwY09HpiT1a5DamE7z57oS6intZpHt3RaLoFefjbuLpvtPdgGeAC+J2Q301YLdXjRuyZxc/7TL0i9XCRdV9L+AABwM5iaIGE3bri9GVCoYC4c1Xn8sY1W5Oki7rMeRtN5Zbb+DvvHRcDpMKOeNTrbZE6BpBFE68jw+LYQnZ5UmelAFxR3MU3zICBtjEUJBpI4F5WQJVZYBikxAvNCsUZ4UFtFvQnO0Mm/uYlhiy3dMXp+Pva0IZjIAdDCuEiI2sB7WfChWFYqR2twwtt7CvBQCzz9gm06GSGR89jWqCfvwvIHuP7+INdPTY3OBRI39I5PLuPYqhBJ2gllEZjQLmebtdKINhUmGuSnC/lL9+wbh66uGd08m2keSIvbEaMt+7keFCLL11AfK3a0Dttm18r1PmMXiJjT7uHIiT7bQr8d7aVCUuZUbuH8/kdIny+psQERJw/4niiPXZbw9feFbHWfABCfXCyGkwmn24OYO1PbpsJfBWsTjfh9Zk+BtET5tFZDPuSx/WNohdnMc+sqxeTxrJrlqK9vXNxJWTW9VlxjRTl4wK2hkDoLHtbxNHkTn8ZiCAEDhRmvnKHCwB8m3xUWcE55uh54bnE65zF4Fvm4GMx5rIZjVrSNjUKGBEjgbejtGrIf3P1VuPmTgTqwsqvlhWROS7fsU6hWyOhuIMMfjqqfrDdJCriA8LG0I2b/I+ENlkSaQlCV/7jrhEDOe7ictxKjcfpmF4CgVfpJ5BGP4OS9uN+dFWHM98j44TvFehwv3Hne9jf8Op8tirHMoIjl0BQGRZwxNMH03OWh0uBGExJTAjBgkqhkiG9w0BCRUxFgQUtwiclhGOgs27XT1wbXZQDj0yyOAwMTAhMAkGBSsOAwIaBQAEFNrsnFXoQilTe1H6GNHplNz6wVIzBAhLV3Iz4VLSkAICCAA=")
		me.CertificateFormat = CertificateFormats.Pkcs12.Ref()
		me.Password = opt.NewString("winona00")
		return []string{credsNotProvided}
	case CredentialsTypes.Token:
		me.Token = opt.NewString("################")
		return []string{credsNotProvided}
	case CredentialsTypes.UsernamePassword:
		me.Username = opt.NewString("################")
		me.Password = opt.NewString("winona00")
		return []string{credsNotProvided}
	}
	return nil
}
