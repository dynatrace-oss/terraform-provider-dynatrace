---
layout: ""
page_title: "dynatrace_azure_connection_authentication Resource - terraform-provider-dynatrace"
subcategory: "Connections"
description: |-
  The resource `dynatrace_azure_connection_authentication` configures federated identity credentials for Azure connections
---

# dynatrace_azure_connection_authentication (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

-> This resource requires the OAuth scopes **Read settings** (`settings:objects:read`) and **Write settings** (`settings:objects:write`)

## Requirements
This resource must be used in combination with the `dynatrace_azure_connection` resource to manage a connection using federated identity credential.
An example of how to set up both resources can be found in the [Resource Example Usage](#resource-example-usage) section below.

## Limitations
If you are creating the `dynatrace_azure_connection_authentication` and the `azuread_application_federated_identity_credential` it is referencing in the same invocation of `terraform apply`, be aware that due to eventual consistency in Azure, the creation of the federated identity credential might not be fully propagated when the `dynatrace_azure_connection_authentication` resource is being created. To mitigate this, we retry the creation of the `dynatrace_azure_connection_authentication` resource with a default timeout of 2 minutes.
If you desire a different timeout, you can adjust it using the `timeouts` block in the `dynatrace_azure_connection_authentication` resource as shown in the example below.
    
The following is an example of the error you might encounter due to this limitation:
```
Error: AADSTS00000: The client '01233456-0123-0123-0123-012345678901'(Application) has no configured federated identity credentials. Trace ID: 01233456-0123-0123-0123-012345678901 Correlation ID: 01233456-0123-0123-0123-012345678901 Timestamp: 2025-12-04 15:34:59Z
```

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

- `terraform-provider-dynatrace -export dynatrace_azure_connection_authentication` downloads all existing Azure connections.

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

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

- `application_id` (String) Application (client) ID of your app registered in Microsoft Azure App registrations
- `azure_connection_id` (String) The ID of a `dynatrace_azure_connection` resource instance for which to define the Azure Authentication
- `directory_id` (String) Directory (tenant) ID of Microsoft Entra ID

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
