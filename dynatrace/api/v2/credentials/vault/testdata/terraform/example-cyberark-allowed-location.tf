// contains the key and the certificate
variable "CREDENTIAL_VAULT_CERT" {
  sensitive = true
}

variable "CREDENTIAL_VAULT_PWD" {
  sensitive = true
}

resource "dynatrace_credentials" "certificate_credentials" {
  name              = "#name#"
  certificate       = base64encode(var.CREDENTIAL_VAULT_CERT)
  format            = "PEM"
  owner_access_only = true
  password          = base64encode(var.CREDENTIAL_VAULT_PWD)
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials" "cyberark_allowed_location" {
  name              = "#name#"
  owner_access_only = true
  external {
    vault_url                       = "https://example.com"
    application_id                  = "my-application-id"
    safe_name                       = "my-safe-name"
    folder_name                     = "my-folder-name"
    account_name                    = "my-account-name"
    certificate                     = dynatrace_credentials.certificate_credentials.id
    location_for_synchronization_id = data.dynatrace_synthetic_location.location.id
  }
  scopes = ["SYNTHETIC"]
}

data "dynatrace_synthetic_location" "location" {
  name = "Location"
}
