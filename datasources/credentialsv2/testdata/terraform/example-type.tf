resource "dynatrace_credentials_v2" "token" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "username_password" {
  name              = "#name#"
  username          = "username"
  password          = "password"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

data "dynatrace_credentials_v2" "by_type" {
  type       = "TOKEN"
  depends_on = [dynatrace_credentials_v2.token, dynatrace_credentials_v2.username_password]
}

# The tenant may hold further TOKEN credentials, so only the ones created here are compared.
output "matched_ids" {
  value = jsonencode(sort(setintersection(
    data.dynatrace_credentials_v2.by_type.items[*].id,
    [dynatrace_credentials_v2.token.id, dynatrace_credentials_v2.username_password.id],
  )))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.token.id])
}
