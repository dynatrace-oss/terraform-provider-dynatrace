resource "dynatrace_openpipeline_spans" "spans" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_spans_#name#"
      data_extraction {
        processor {
          bizevent_extraction_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_Business_event_processor_#name#"
            matcher     = "true"
            event_provider {
              type  = "field"
              field = "Provider"
            }
            event_type {
              type  = "field"
              field = "Something"
            }
            field_extraction {
              semantic = "INCLUDE_ALL"
            }
          }
        }
      }
      processing {
        processor {
          fields_add_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_Add_field_#name#"
            matcher     = "true"
            field {
              name  = "test"
              value = "1"
            }
          }
        }
      }
    }
  }
}
