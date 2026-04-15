resource "dynatrace_openpipeline_v2_security_events_pipelines" "pipeline" {
  display_name = "Warning pipeline"
  custom_id = "pipeline_Warning_pipeline_2773_tf_#name#"
  metric_extraction {
    processors {
      processor {
        type        = "counterMetric"
        id          = "processor_Count_warning_events_6392"
        description = "Count warnings"
        matcher     = "true"
        counter_metric {
          metric_key = "warning.count"
          dimensions {
            dimension {
              source_field_name = "dt.cost.costcenter"
            }
            # This dimension is updated in this test case
            dimension {
              source_field_name = "dt.cost.product.edit"
            }
          }
        }
        enabled = true
      }
    }
  }
}
