resource "dynatrace_autotag_rules" "rules" {
  auto_tag_id = dynatrace_autotag_v2.tag.id
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


resource "dynatrace_autotag_v2" "tag" {
  name                        = "#name#"
  rules_maintained_externally = true
}
