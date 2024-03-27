resource "dynatrace_infraops_app_feature_flags" "#name#" {
  name         = "Terraform"
  type         = "string"
  string_value = "#name#"
}