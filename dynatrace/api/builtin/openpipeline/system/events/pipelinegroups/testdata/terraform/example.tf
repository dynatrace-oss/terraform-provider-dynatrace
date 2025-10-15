resource "dynatrace_openpipeline_v2_system_events_pipelines" "example1" {
  display_name = "#name#"
  custom_id = "#name#"
  processing {}
  davis {
    processors {
      processor {
        type        = "davis"
        id          = "processor_Create_warning_event_8226"
        description = "Create warning event"
        matcher     = "true"
        davis {
          properties {
            property {
              key = "event.type"
              value = "CUSTOM_ALERT"
            }
            property {
              key = "event.name"
              value = "Warning detected"
            }
            property {
              key = "event.description"
              value = "Warning: {dims:record.summary}"
            }
          }
        }
        enabled = true
      }
    }
  }
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
              source_field_name = "dt.cost.product"
            }
            dimension {
              source_field_name = "dt.security_context"
            }
            dimension {
              source_field_name      = "record.category"
              destination_field_name = "warning_category"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "valueMetric"
        id          = "processor_Warning_timeout_1990"
        description = "Warning timeout"
        matcher     = "true"
        value_metric {
          metric_key    = "warning.timeout"
          field         = "recording.timeout_in_min"
          default_value = 60
          dimensions {
            dimension {
              source_field_name = "dt.cost.costcenter"
            }
            dimension {
              source_field_name = "dt.cost.product"
            }
            dimension {
              source_field_name = "dt.security_context"
            }
            dimension {
              source_field_name      = "record.category"
              destination_field_name = "warning_category"
            }
          }
        }
        enabled = true
      }
    }
  }
  security_context {}
  cost_allocation {}
  product_allocation {}
  storage {}
  data_extraction {}
}


resource "dynatrace_openpipeline_v2_system_events_pipelines" "example2" {
  display_name = "#name#-2"
  custom_id = "#name#-2"
  processing {}
  product_allocation {}
  data_extraction {}
  metric_extraction {}
  security_context {}
  storage {}
  cost_allocation {}
  davis {}
}

resource "dynatrace_openpipeline_v2_system_events_pipelinegroups" "example" {
  display_name = "#name#"
  included_pipelines {
    included_pipeline {
      is_target_pipeline_placeholder = false
      pipeline_id = dynatrace_openpipeline_v2_system_events_pipelines.example1.id
      enabled_stages = ["davis", "metricExtraction"]
    }
  }
  targeted_pipelines = [dynatrace_openpipeline_v2_system_events_pipelines.example2.id]
}
