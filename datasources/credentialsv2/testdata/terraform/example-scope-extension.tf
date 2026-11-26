// contains the key and the certificate
variable "CREDENTIAL_VAULT_CERT" {
  sensitive = true
}

variable "CREDENTIAL_VAULT_PWD" {
  sensitive = true
}

resource "dynatrace_credentials_v2" "synthetic" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "extension" {
  name              = "#name#"
  certificate       = base64encode(var.CREDENTIAL_VAULT_CERT)
  password          = base64encode(var.CREDENTIAL_VAULT_PWD)
  format            = "PEM"
  public            = true
  owner_access_only = false
  scopes            = ["EXTENSION"]
}

data "dynatrace_credentials_v2" "by_scope" {
  scope      = "EXTENSION"
  depends_on = [dynatrace_credentials_v2.synthetic, dynatrace_credentials_v2.extension]
}

# The tenant may hold further EXTENSION credentials, so only the ones created here are compared.
output "matched_ids" {
  value = jsonencode(sort(setintersection(
    data.dynatrace_credentials_v2.by_scope.items[*].id,
    [dynatrace_credentials_v2.synthetic.id, dynatrace_credentials_v2.extension.id],
  )))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.extension.id])
}
