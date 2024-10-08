resource "dynatrace_openpipeline_business_events" "bizevents" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "#name#"
      processing {
        processor {
          fields_remove_processor {
            description = "#name#"
            enabled     = true
            fields      = ["#name#"]
            id          = "#name#"
            matcher     = "true"
          }
        }
      }
    }
  }
}