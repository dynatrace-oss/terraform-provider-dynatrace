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
  name = "Example Azure Connector"
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
