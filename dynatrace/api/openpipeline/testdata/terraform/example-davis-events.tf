resource "dynatrace_openpipeline_davis_events" "davis_events" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_davis_events_#name#"
      storage {
        catch_all_bucket_name = "default_davis_custom_events"
        processor {
          no_storage_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_No_storage_#name#"
            matcher     = "true"
          }
        }
      }
    }
  }
}