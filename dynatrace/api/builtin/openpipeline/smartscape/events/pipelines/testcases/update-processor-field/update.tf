resource "dynatrace_openpipeline_v2_smartscape_events_pipelines" "pipeline" {
  display_name = "Warning pipeline"
  custom_id    = "pipeline_Warning_pipeline_2773_tf_#name#"
  processing {
    processors {
      processor {
        type        = "fieldsAdd"
        id          = "processor_Add_warning_flag_5434"
        description = "Add warning flag"
        matcher     = "true"
        fields_add {
          fields {
            field {
              name       = "custom.is_warning"
              value      = "true"
              value_type = "boolean"
            }
            # This field is updated in this test case
            field {
              name       = "custom.category"
              value      = "category-edit"
              value_type = "string"
            }
          }
        }
        enabled = true
      }
    }
  }
}
