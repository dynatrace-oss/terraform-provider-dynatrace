resource "dynatrace_log_custom_source" "source" {
  name    = "#name#"
  enabled = false
  scope   = "HOST_GROUP-0000000000000000"
  context {
    # update => re-create due to set-hash change
    context {
      attribute = "dt.entity.process_group"
      values = ["PROCESS_GROUP-0000000000000001"]
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
      # update => re-create due to set-hash change
      custom_log_source_with_enrichment {
        path = "/terraformEdit"
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
          # update => re-create due to set-hash change
          enrichment {
            type  = "attribute"
            key   = "keyEdit"
            value = "valueEdit"
          }
        }
      }
    }
  }
}
