// contains the key and the certificate
variable "CREDENTIAL_VAULT_CERT" {
  sensitive = true
}

variable "CREDENTIAL_VAULT_PWD" {
  sensitive = true
}

resource "dynatrace_credentials_v2" "certificate" {
  name              = "#name#"
  certificate       = base64encode(var.CREDENTIAL_VAULT_CERT)
  password          = base64encode(var.CREDENTIAL_VAULT_PWD)
  format            = "PEM"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "public_certificate" {
  name              = "#name#"
  certificate       = base64encode(var.CREDENTIAL_VAULT_CERT)
  password          = base64encode(var.CREDENTIAL_VAULT_PWD)
  format            = "PEM"
  public            = true
  owner_access_only = false
  scopes            = ["EXTENSION"]
}

data "dynatrace_credentials_v2" "by_type" {
  type       = "PUBLIC_CERTIFICATE"
  depends_on = [dynatrace_credentials_v2.certificate, dynatrace_credentials_v2.public_certificate]
}

# The tenant may hold further PUBLIC_CERTIFICATE credentials, so only the ones created here are compared.
output "matched_ids" {
  value = jsonencode(sort(setintersection(
    data.dynatrace_credentials_v2.by_type.items[*].id,
    [dynatrace_credentials_v2.certificate.id, dynatrace_credentials_v2.public_certificate.id],
  )))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.public_certificate.id])
}
