resource "dynatrace_log_custom_source" "source" {
  name    = "#name#"
  enabled = false
  scope   = "HOST_GROUP-0000000000000000"
  context {
    context {
      attribute = "dt.entity.process_group"
      values = ["PROCESS_GROUP-0000000000000000"]
    }
  }
  custom_log_source {
    type = "WINDOWS_EVENT_LOG"
    values_and_enrichment {
      custom_log_source_with_enrichment {
        path = "/terraform1"
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
      custom_log_source_with_enrichment {
        path = "/terraform2"
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
      custom_log_source_with_enrichment {
        path = "/terraform3"
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
