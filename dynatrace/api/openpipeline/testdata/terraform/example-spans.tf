resource "dynatrace_openpipeline_spans" "spans" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "#name#"
      id           = "pipeline_Custom_spans_#name#"
      data_extraction {
        processor {
          bizevent_extraction_processor {
            description = "Custom bizevent extraction"
            enabled     = true
            id          = "processor_custom_bizevent_1_#name#"
            matcher     = "true"
            sample_data = "{}"
            field_extraction {
              semantic = "INCLUDE"
              fields = ["my.field"]
            }
            event_provider {
              type = "constant"
              constant = "my-constant"
            }
            event_type {
              type = "constant"
              constant = "my-constant"
            }
          }
        }
      }
      metric_extraction {
        processor {
          sampling_aware_counter_metric_extraction_processor {
            description = "Custom sampling counter extraction"
            enabled     = true
            id          = "processor_custom_sampling_counter_1_#name#"
            matcher     = "true"
            metric_key  = "events.counter"
            aggregation = "ENABLED"
            sample_data = "{}"
            sampling = "ENABLED"
            dimensions = ["ab=xy"]
          }
        }
        processor {
          sampling_aware_value_metric_extraction_processor {
            description = "Custom sampling value extraction"
            enabled     = true
            id          = "processor_custom_sampling_value_1_#name#"
            matcher     = "true"
            measurement = "FIELD"
            metric_key  = "events.value"
            aggregation = "DISABLED"
            sampling = "DISABLED"
            default_value = "10"
            field = "my.field"
            sample_data = "{}"
            dimensions = ["xyz=abc"]
          }
        }
      }
      processing {
        processor {
          fields_add_processor {
            description = "#name#"
            enabled     = true
            id          = "processor_Add_field_#name#"
            matcher     = "true"
            field {
              name  = "test"
              value = "1"
            }
          }
        }
      }
    }
  }
}
