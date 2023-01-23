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

// CertificateFormat The certificate format.
type CertificateFormat string

// CertificateFormats offers the known enum values
var CertificateFormats = struct {
	Pem     CertificateFormat
	Pkcs12  CertificateFormat
	Unknown CertificateFormat
}{
	"PEM",
	"PKCS12",
	"UNKNOWN",
}

func (me CertificateFormat) Ref() *CertificateFormat {
	return &me
}

// CredentialsResponseElementType The type of the credentials set.
type CredentialsResponseElementType string

// CredentialsResponseElementTypes offers the known enum values
var CredentialsResponseElementTypes = struct {
	Certificate       CredentialsResponseElementType
	PublicCertificate CredentialsResponseElementType
	Token             CredentialsResponseElementType
	Unknown           CredentialsResponseElementType
	UsernamePassword  CredentialsResponseElementType
}{
	CredentialsResponseElementType("CERTIFICATE"),
	CredentialsResponseElementType("PUBLIC_CERTIFICATE"),
	CredentialsResponseElementType("TOKEN"),
	CredentialsResponseElementType("UNKNOWN"),
	CredentialsResponseElementType("USERNAME_PASSWORD"),
}

type MonitorType string

var MonitorTypes = struct {
	HTTPMonitor     MonitorType
	BrowserMonitorl MonitorType
}{
	MonitorType("HTTP_MONITOR"),
	MonitorType("BROWSER_MONITOR"),
}

type ExternalVaultConfigType string

var ExternalVaultConfigTypes = struct {
	AzureCertificateModel    ExternalVaultConfigType
	AzureClientSecretModel   ExternalVaultConfigType
	HashcorpAppRoleModel     ExternalVaultConfigType
	HashcorpCertificateModel ExternalVaultConfigType
}{
	ExternalVaultConfigType("AZURE_CERTIFICATE_MODEL"),
	ExternalVaultConfigType("AZURE_CLIENT_SECRET_MODEL"),
	ExternalVaultConfigType("HASHICORP_APPROLE_MODEL"),
	ExternalVaultConfigType("HASHICORP_CERTIFICATE_MODEL"),
}

type CredentialsType string

var CredentialsTypes = struct {
	Certificate       CredentialsType
	PublicCertificate CredentialsType
	Token             CredentialsType
	UsernamePassword  CredentialsType
	Unknown           CredentialsType
}{
	"CERTIFICATE",
	"PUBLIC_CERTIFICATE",
	"TOKEN",
	"USERNAME_PASSWORD",
	"UNKNOWN",
}

type Scope string

var Scopes = struct {
	All       Scope
	Extension Scope
	Synthetic Scope
	Unknown   Scope
}{
	Scope("ALL"),
	Scope("EXTENSION"),
	Scope("SYNTHETIC"),
	Scope("UNKNOWN"),
}
