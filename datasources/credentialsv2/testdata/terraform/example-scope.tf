resource "dynatrace_credentials_v2" "synthetic" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["SYNTHETIC"]
}

resource "dynatrace_credentials_v2" "app_engine" {
  name              = "#name#"
  token             = "token"
  owner_access_only = false
  scopes            = ["APP_ENGINE"]
}

data "dynatrace_credentials_v2" "by_scope" {
  scope      = "APP_ENGINE"
  depends_on = [dynatrace_credentials_v2.synthetic, dynatrace_credentials_v2.app_engine]
}

# The tenant may hold further APP_ENGINE credentials, so only the ones created here are compared.
output "matched_ids" {
  value = jsonencode(sort(setintersection(
    data.dynatrace_credentials_v2.by_scope.items[*].id,
    [dynatrace_credentials_v2.synthetic.id, dynatrace_credentials_v2.app_engine.id],
  )))
}

output "expected_ids" {
  value = jsonencode([dynatrace_credentials_v2.app_engine.id])
}
