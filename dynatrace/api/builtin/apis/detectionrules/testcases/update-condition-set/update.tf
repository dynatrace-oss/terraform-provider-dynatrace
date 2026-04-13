resource "dynatrace_api_detection" "detection" {
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
    condition {
      base    = "PACKAGE"
      matcher = "BEGINS_WITH"
      pattern = "com.terraform3"
    }
    condition {
      base    = "PACKAGE"
      matcher = "BEGINS_WITH"
      pattern = "com.terraformNew"
    }
  }
}
