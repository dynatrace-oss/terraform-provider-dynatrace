resource "dynatrace_api_detection" "#name#" {
  api_color       = "#5ead35"
  api_name        = "#name#"
  technology      = "Go"
  third_party_api = false
  conditions {
    condition {
      base    = "PACKAGE"
      matcher = "BEGINS_WITH"
      pattern = "com.terraform"
    }
  }
}