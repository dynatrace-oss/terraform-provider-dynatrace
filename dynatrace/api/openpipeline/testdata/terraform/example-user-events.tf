resource "dynatrace_openpipeline_user_events" "user_events" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_user_events_#name#"
      security_context {
        processor {
          security_context_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_Set_security_context_#name#"
            matcher     = "true"
            value {
              type     = "constant"
              constant = "val"
            }
          }
        }
      }
    }
  }
}
