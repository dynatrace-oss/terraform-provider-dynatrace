resource "dynatrace_log_custom_source" "#name#" {
  name    = "#name#"
  enabled = false
  scope   = "HOST_GROUP-1234567890000000"
  custom_log_source {
    type = "WINDOWS_EVENT_LOG"
    values_and_enrichment {
      custom_log_source_with_enrichment {
        path = "/terraform"
        enrichment {
          enrichment {
            type  = "attribute"
            key   = "key1"
            value = "value1"
          }
          enrichment {
            type  = "attribute"
            key   = "key2"
            value = "value2"
          }
        }
      }
    }
  }
}
