resource "dynatrace_autotag_rules" "Avengers" {
  auto_tag_id = "vu9U3hXa3q0AAAABABlidWlsdGluOnRhZ3MuYXV0by10YWdnaW5nAAZ0ZW5hbnQABnRlbmFudAAkOTNjMDBkYzktY2FkMC0zMWY3LWEzZGQtMWQ4MDdjMWQwMjhivu9U3hXa3q0"
  rules {
    rule {
      type                = "ME"
      enabled             = true
      value_format        = "Avengers"
      value_normalization = "Leave text as-is"
      attribute_rule {
        entity_type                 = "SERVICE"
        service_to_host_propagation = true
        service_to_pgpropagation    = true
        conditions {
          condition {
            case_sensitive = true
            key            = "SERVICE_DATABASE_NAME"
            operator       = "EQUALS"
            string_value   = "AvengersA-1"
          }
          condition {
            case_sensitive = true
            key            = "SERVICE_DATABASE_NAME"
            operator       = "EQUALS"
            string_value   = "AvengersA-2"
          }
        }
      }
    }
    rule {
      type                = "ME"
      enabled             = true
      value_format        = "Avengers"
      value_normalization = "Leave text as-is"
      attribute_rule {
        entity_type                 = "SERVICE"
        service_to_host_propagation = true
        service_to_pgpropagation    = true
        conditions {
          condition {
            case_sensitive = true
            key            = "SERVICE_DATABASE_NAME"
            operator       = "EQUALS"
            string_value   = "AvengersB-1"
          }
          condition {
            case_sensitive = true
            key            = "SERVICE_DATABASE_NAME"
            operator       = "EQUALS"
            string_value   = "AvengersB-2"
          }
          condition {
            case_sensitive = true
            key            = "SERVICE_DATABASE_NAME"
            operator       = "EQUALS"
            string_value   = "AvengersB-3"
          }
        }
      }
    }
  }
}
