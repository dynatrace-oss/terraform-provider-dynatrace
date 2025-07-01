resource "dynatrace_openpipeline_davis_problems" "davis_problems" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_davis_problems_#name#"
      processing {
        processor {
          fields_rename_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_Rename_problem_ID_#name#"
            matcher     = "true"
            field {
              from_name = "problem_id"
              to_name   = "problemId"
            }
          }
        }
      }
    }
  }
}
