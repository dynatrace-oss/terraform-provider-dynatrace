variable "CUSTOM_ID" {
  type = string
}
# contains removed and added dimensions
resource "dynatrace_openpipeline_v2_bizevents_pipelines" "max-pipeline" {
  display_name = "Warning pipeline"
  custom_id    = var.CUSTOM_ID
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
            dimension {
              source_field_name = "dt.cost.product.new"
            }
          }
        }
        enabled = true
      }
    }
  }
}
