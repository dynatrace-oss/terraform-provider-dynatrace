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

resource "dynatrace_credentials_v2" "hashicorp_certificate" {
  name              = "#name#"
  owner_access_only = true
  external {
    vault_url           = "https://example.com:8020"
    path_to_credentials = "secret/my-token"
    certificate         = dynatrace_credentials_v2.certificate_credentials.id
    token_secret_name   = "my-token-secret"
  }
  scopes = ["SYNTHETIC"]
}
