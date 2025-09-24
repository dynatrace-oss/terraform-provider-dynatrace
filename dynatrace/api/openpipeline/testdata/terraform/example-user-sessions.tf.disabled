resource "dynatrace_openpipeline_user_sessions" "usersessions" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_user_sessions_#name#"
      security_context {
        processor {
          security_context_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_Set_to_static_#name#"
            matcher     = "true"
            value {
              type     = "constant"
              constant = "test"
            }
          }
        }
      }
    }
  }
}