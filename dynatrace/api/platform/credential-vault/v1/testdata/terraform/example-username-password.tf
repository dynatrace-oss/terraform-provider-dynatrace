resource "dynatrace_credentials_v2" "username_password_credentials" {
  name                       = "#name#"
  username                   = "username"
  password                   = "password"
  owner_access_only          = true
  scopes                     = ["APP_ENGINE", "SYNTHETIC"]
  allow_contextless_requests = true
}
