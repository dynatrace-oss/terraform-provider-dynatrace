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
  scopes            = ["APP_ENGINE"]
}

data "dynatrace_credentials_v2" "all" {
  depends_on = [dynatrace_credentials_v2.token, dynatrace_credentials_v2.username_password]
}

# The tenant may hold further credentials, so only the ones created here are compared.
output "matched_ids" {
  value = jsonencode(sort(setintersection(
    data.dynatrace_credentials_v2.all.items[*].id,
    [dynatrace_credentials_v2.token.id, dynatrace_credentials_v2.username_password.id],
  )))
}

output "expected_ids" {
  value = jsonencode(sort([dynatrace_credentials_v2.token.id, dynatrace_credentials_v2.username_password.id]))
}
