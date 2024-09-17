resource "dynatrace_request_naming" "team-hawaiian-pizza" {
  enabled        = false
  naming_pattern = "team-hawaiian-pizza"
  conditions {
    condition {
      attribute = "WEBREQUEST_URL"
      comparison {
        string {
          operator = "EQUALS"
          value    = "team-hawaiian-pizza"
        }
      }
    }
  }
}
