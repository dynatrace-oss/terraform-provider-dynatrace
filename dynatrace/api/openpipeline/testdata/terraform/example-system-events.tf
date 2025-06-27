resource "dynatrace_openpipeline_system_events" "system_events" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_system_events_#name#"
      data_extraction {
        processor {
          davis_event_extraction_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_My_Davis_event_#name#"
            matcher     = "true"
            properties {
              key   = "event.type"
              value = "CUSTOM_ALERT"
            }
            properties {
              key   = "event.name"
              value = "test.event"
            }
            properties {
              key   = "var"
              value = "val"
            }
            properties {
              key   = "event.description"
              value = "Some description"
            }
          }
        }
      }
    }
  }
}
