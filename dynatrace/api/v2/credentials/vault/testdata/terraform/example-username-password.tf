resource "dynatrace_credentials" "username_password_credentials" {
  name              = "#name#"
  username          = "username"
  password          = "password"
  owner_access_only = true
  scopes            = ["SYNTHETIC"]
}
