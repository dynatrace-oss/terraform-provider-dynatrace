resource "dynatrace_openpipeline_logs" "logs" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "test"
      id           = "pipeline_test_logs_#name#"
      processing {
        processor {
          fields_rename_processor {
            description = "test"
            enabled     = true
            id          = "processor_test_logs_#name#"
            matcher     = "true"
            field {
              from_name = "#name#"
              to_name   = "#name#bar"
            }
          }
        }
      }
    }
  }
}