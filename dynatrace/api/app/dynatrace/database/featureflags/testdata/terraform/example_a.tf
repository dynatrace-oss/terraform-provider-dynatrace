resource "dynatrace_db_app_feature_flags" "#name#" {
  name         = "Terraform"
  type         = "string"
  string_value = "#name#"
}