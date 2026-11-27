resource "dynatrace_credentials_v2" "token_synthetic" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "token_app_engine" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["APP_ENGINE"]
}

resource "dynatrace_credentials_v2" "username_password_synthetic" {
  name              = "#name#"
  username          = "username"
  password          = "password"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

data "dynatrace_credentials_v2" "by_type_and_scope" {
  type  = "TOKEN"
  scope = "SYNTHETIC"
  depends_on = [
    dynatrace_credentials_v2.token_synthetic,
    dynatrace_credentials_v2.token_app_engine,
    dynatrace_credentials_v2.username_password_synthetic,
  ]
}

# The tenant may hold further matching credentials, so only the ones created here are compared.
output "matched_ids" {
  value = jsonencode(sort(setintersection(
    data.dynatrace_credentials_v2.by_type_and_scope.items[*].id,
    [
      dynatrace_credentials_v2.token_synthetic.id,
      dynatrace_credentials_v2.token_app_engine.id,
      dynatrace_credentials_v2.username_password_synthetic.id,
    ],
  )))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.token_synthetic.id])
}
