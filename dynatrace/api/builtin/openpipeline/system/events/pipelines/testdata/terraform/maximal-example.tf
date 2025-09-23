resource "dynatrace_openpipeline_v2_system_events_pipelines" "max-pipeline" {
  display_name = "Warning pipeline"
  custom_id = "pipeline_Warning_pipeline_2773_tf_#name#"
  processing {
    # processing is not available for system_events pipelines
  }
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
  security_context {
    # security_context is not available for system_events pipelines
  }
  cost_allocation {

  }
  product_allocation {

  }
  storage {
    # storage is not supported for system_events pipelines
  }
  data_extraction {}
}
