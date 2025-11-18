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
  name = "Example Azure Connector - Secret"
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
