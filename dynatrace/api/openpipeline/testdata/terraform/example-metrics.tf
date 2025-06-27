resource "dynatrace_openpipeline_metrics" "metrics" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_metrics_#name#"
      processing {
        processor {
          drop_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_Drop_all_records_#name#"
            matcher     = "true"
          }
        }
      }
    }
  }
}