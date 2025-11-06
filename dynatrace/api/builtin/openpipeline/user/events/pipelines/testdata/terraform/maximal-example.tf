resource "dynatrace_openpipeline_v2_user_events_pipelines" "max-pipeline" {
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
    processors {
      processor {
        type        = "securityContext"
        id          = "processor_Use_dt.security_context_if_set_1080"
        description = "Use dt.security_context if set"
        matcher     = "isNotNull(dt.security_context)"
        security_context {
          value {
            type = "field"
            field {
              source_field_name = "dt.security_context"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "securityContext"
        id          = "processor_Assign_warnings_to_ACME_teams_if_no_context_set_5465"
        description = "Assign warnings to ACME teams if no context set"
        matcher     = "isNull(dt.security_context)"
        security_context {
          value {
            type = "multiValueConstant"
            multi_value_constant = ["ACME1", "ACME2"]
          }
        }
        enabled = true
      }
    }
  }
}
