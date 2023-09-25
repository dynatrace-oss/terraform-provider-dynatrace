resource "dynatrace_platform_bucket" "#name#" {
  name         = "#name#"
  display_name = "Custom logs bucket playground"
  retention    = 67
  table        = "spans"
}
