// contains the key and the certificate
variable "CREDENTIAL_VAULT_CERT" {
  sensitive = true
}

variable "CREDENTIAL_VAULT_PWD" {
  sensitive = true
}

resource "dynatrace_credentials_v2" "certificate_credentials" {
  name              = "#name#"
  certificate       = base64encode(var.CREDENTIAL_VAULT_CERT)
  format            = "PEM"
  owner_access_only = true
  password          = base64encode(var.CREDENTIAL_VAULT_PWD)
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "cyberark_allowed_location" {
  name              = "#name#"
  owner_access_only = true
  external {
    vault_url      = "https://example.com"
    application_id = "my-application-id"
    safe_name      = "my-safe-name"
    folder_name    = "my-folder-name"
    account_name   = "my-account-name"
    certificate    = dynatrace_credentials_v2.certificate_credentials.id
  }
  scopes = ["SYNTHETIC"]
}
