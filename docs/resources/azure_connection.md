---
layout: ""
page_title: "dynatrace_azure_connection Resource - terraform-provider-dynatrace"
subcategory: "Connections"
description: |-
  The resource `dynatrace_azure_connection` covers configuration for Azure connections
---

# dynatrace_azure_connection (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

-> This resource requires the OAuth scopes **Read settings** (`settings:objects:read`) and **Write settings** (`settings:objects:write`)

## Requirements
This resource can be used to create connections using an Azure client secret or federated identity credential. For the latter case, this resource must be used together with a `dynatrace_azure_connection_authentication` resource.
Ensure you configure both resources together for a valid Azure connection.
An example of how to set up a connection using a client secret or a federated identiy credential can be found in the [Resource Example Usage](#resource-example-usage) section below.

## Limitations
~> **Warning** If a resource is created using an API token or without setting `DYNATRACE_HTTP_OAUTH_PREFERENCE=true` (when both are used), the settings object's owner will remain empty.

An empty owner implies:
- The settings object becomes public, allowing other users with settings permissions to read and modify it.
- Changing the settings object's permissions will have no effect, meaning the `dynatrace_settings_permissions` resource can't alter its access.

When a settings object is created using platform credentials:
- The owner is set to the owner of the OAuth client or platform token.
- By default, the settings object is private; only the owner can read and modify it.
- Access modifiers can be managed using the `dynatrace_settings_permissions` resource.

We recommend using platform credentials to ensure a correct setup.
In case an API token is needed, we recommend setting `DYNATRACE_HTTP_OAUTH_PREFERENCE=true`.

## Dynatrace Documentation

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:hyperscaler-authentication.connections.azure`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_azure_connection` downloads all existing Azure connections.

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

### Azure connection with client secret authentication
```terraform
variable "azure_tenant_id" {
  type        = string
  description = "The Azure Active Directory tenant ID."
}

terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 3.1.0"
    }
  }
}

# Configure the Azure Active Directory provider
provider "azuread" {
  tenant_id = var.azure_tenant_id
}

# Create an application
resource "azuread_application_registration" "example" {
  display_name = "ExampleApp"
}

# Create a client secret
resource "azuread_application_password" "example" {
  application_id = azuread_application_registration.example.id
}

# Create Azure connection
resource "dynatrace_azure_connection" "example" {
  name = "#name#"
  type = "clientSecret"
  client_secret {
    client_secret  = azuread_application_password.example.value
    application_id = azuread_application_registration.example.client_id
    directory_id   = var.azure_tenant_id
    consumers = [
      "DA"
    ]
  }
}
```

### Azure connection with federated identity credential authentication
```terraform
variable "azure_tenant_id" {
  type        = string
  description = "The Azure Active Directory tenant ID."
}

variable "dynatrace_environment_url" {
  type        = string
  description = "The Dynatrace environment URL"
}

variable "dynatrace_token_issuer" {
  type        = string
  description = "The Dynatrace token issuer URL"
}


terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 3.1.0"
    }
  }
}

# Configure the Azure Active Directory provider
provider "azuread" {
  tenant_id = var.azure_tenant_id
}


# Create an application
resource "azuread_application_registration" "example" {
  display_name = "ExampleApp"
}

# Create basic Azure connection
resource "dynatrace_azure_connection" "example" {
  name = "#name#"
  type = "federatedIdentityCredential"
  federated_identity_credential {
    consumers = [
      "APP:dynatrace.microsoft.azure.connector"
    ]
  }
}

# Create a federated identity credential
resource "azuread_application_federated_identity_credential" "example" {
  application_id = azuread_application_registration.example.id
  display_name   = "Example"
  audiences      = ["${var.dynatrace_environment_url}/app-id/dynatrace.microsoft.azure.connector"]
  issuer         = var.dynatrace_token_issuer
  subject        = "dt:connection-id/${dynatrace_azure_connection.example.id}"
}

# Update the Azure connection with authentication details
resource "dynatrace_azure_connection_authentication" "example" {
  azure_connection_id = dynatrace_azure_connection.example.id
  application_id      = azuread_application_registration.example.client_id
  directory_id        = var.azure_tenant_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the connection
- `type` (String) Azure Authentication mechanism to be used by the connection. Possible Values: `clientSecret`, `federatedIdentityCredential`

### Optional

- `client_secret` (Block List, Max: 1) no documentation available (see [below for nested schema](#nestedblock--client_secret))
- `federated_identity_credential` (Block List, Max: 1) no documentation available (see [below for nested schema](#nestedblock--federated_identity_credential))
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--client_secret"></a>
### Nested Schema for `client_secret`

Required:

- `application_id` (String) Application (client) ID of your app registered in Microsoft Azure App registrations
- `client_secret` (String, Sensitive) Client secret of your app registered in Microsoft Azure App registrations
- `directory_id` (String) Directory (tenant) ID of Microsoft Entra ID

Optional:

- `consumers` (List of String) Dynatrace integrations that can use this connection. Possible Values: `DA`, `NONE`, `SVC:com.dynatrace.da`


<a id="nestedblock--federated_identity_credential"></a>
### Nested Schema for `federated_identity_credential`

Optional:

- `consumers` (List of String) Consumers that can use the connection. Possible Values: `APP:dynatrace.microsoft.azure.connector`, `DA`, `NONE`, `SVC:com.dynatrace.da`, `SVC:com.dynatrace.openpipeline`


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `update` (String)
