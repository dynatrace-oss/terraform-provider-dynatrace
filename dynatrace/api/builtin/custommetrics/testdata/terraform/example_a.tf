resource "dynatrace_user_session_metrics" "metric" {
  enabled    = false
  metric_key = "uscm.#name#"
  filters {
    filter {
      field_name = "useraction.application"
      operator   = "EQUALS"
      value      = "www.terraform.io/"
    }
    filter {
      field_name = "useraction.name"
      operator   = "EQUALS"
      value      = "Loading of page /"
    }
  }
  value {
    type = "COUNTER"
  }
}