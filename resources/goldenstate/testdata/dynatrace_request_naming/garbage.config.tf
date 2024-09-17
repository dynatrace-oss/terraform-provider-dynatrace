resource "dynatrace_request_naming" "team-hawaiian-deleteme" {
  enabled        = false
  naming_pattern = "team-hawaiian-deleteme"
  conditions {
    condition {
      attribute = "WEBREQUEST_URL"
      comparison {
        string {
          operator = "EQUALS"
          value    = "team-hawaiian-deleteme"
        }
      }
    }
  }
}
