// contains the key and the certificate
variable "CREDENTIAL_VAULT_CERT" {
  sensitive = true
}

variable "CREDENTIAL_VAULT_PWD" {
  sensitive = true
}


resource "dynatrace_credentials_v2" "public_certificate" {
  name              = "#name#"
  description       = "Root certificate for validating Extension 2.0 signatures"
  certificate       = base64encode(var.CREDENTIAL_VAULT_CERT)
  password          = base64encode(var.CREDENTIAL_VAULT_PWD)
  format            = "PEM"
  owner_access_only = true
  public            = true
  scopes            = ["EXTENSION"]
}
