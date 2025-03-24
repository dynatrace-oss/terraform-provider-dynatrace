resource "dynatrace_service_naming_order" "this" {
  naming_rule_ids = [
    dynatrace_service_naming.first.id,
    dynatrace_service_naming.second.id,
  ]  
}

resource "dynatrace_service_naming" "first" {
  name    = "${randomize}"
  enabled = true
  format  = "${randomize}"
  conditions {
    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_TYPE"
      }
      service_type {
        negate   = false
        operator = "EQUALS"
        value    = "WEB_REQUEST_SERVICE"
      }
    }
  }
}

resource "dynatrace_service_naming" "second" {
  name    = "${randomize}"
  enabled = true
  format  = "${randomize}"
  conditions {
    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_TYPE"
      }
      service_type {
        negate   = false
        operator = "EQUALS"
        value    = "WEB_REQUEST_SERVICE"
      }
    }
  }
}
