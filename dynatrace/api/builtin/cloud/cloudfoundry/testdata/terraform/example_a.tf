resource "dynatrace_cloud_foundry" "#name#" {
  enabled           = true
  active_gate_group = "Terraform"
  api_url           = "https://www.google.at/test/#name#"
  label             = "#name#"
  login_url         = "https://www.google.at/test/#name#"
  password          = "pass2"
  username          = "user"
}