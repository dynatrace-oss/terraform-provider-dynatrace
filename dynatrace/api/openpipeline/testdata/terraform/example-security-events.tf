resource "dynatrace_openpipeline_security_events" "events_security" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_test_security_#name#"
      processing {
        processor {
          fields_remove_processor {
            description = "#name#"
            enabled     = true
            fields      = ["test"]
            id          = "processor_test_security_#name#"
            matcher     = "true"
          }
        }
      }
    }
  }
}