resource "dynatrace_api_detection" "first-entry" {
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

resource "dynatrace_api_detection" "second-entry" {
  api_color       = "#5ead35"
  api_name        = "#name#-second"
  technology      = "Go"
  third_party_api = false
  conditions {
    condition {
      base    = "PACKAGE"
      matcher = "BEGINS_WITH"
      pattern = "com.terraform-second"
    }
  }
  insert_after = "${dynatrace_api_detection.first-entry.id}"
}