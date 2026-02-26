# certificate encoded in base64
variable "certificate" {
  sensitive = true
}

variable "certificate_password" {
  sensitive = true
}

resource "dynatrace_credentials" "certificate_credentials" {
  name              = "#name#"
  certificate       = var.certificate
  format            = "PKCS12"
  owner_access_only = true
  password          = var.certificate_password
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials" "cyberark_allowed_location" {
  name              = "#name#"
  owner_access_only = true
  external {
    vault_url      = "https://example.com"
    application_id = "my-application-id"
    safe_name      = "my-safe-name"
    folder_name    = "my-folder-name"
    account_name   = "my-account-name"
    certificate    = dynatrace_credentials.certificate_credentials.id
  }
  scopes = ["SYNTHETIC"]
}
