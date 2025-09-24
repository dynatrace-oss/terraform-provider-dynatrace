resource "dynatrace_openpipeline_logs" "logs" {
  pipelines {
    pipeline {
      enabled      = true
      display_name = "test"
      id           = "pipeline_test_logs_#name#"
      cost_allocation {
        processor {
          cost_allocation_processor {
            description = "Custom cost allocation 1"
            enabled     = true
            id          = "processor_custom_cost_allocation_1_#name#"
            matcher     = "true"
            sample_data = "{}"
            value {
              type     = "constant"
              constant = "string"
            }
          }
        }
      }
      product_allocation {
        processor {
          product_allocation_processor {
            description = "Custom product allocation 1"
            enabled     = true
            id          = "processor_custom_product_allocation_1_#name#"
            matcher     = "true"
            sample_data = "{}"
            value {
              type     = "constant"
              constant = "string"
            }
          }
        }
      }
      processing {
        processor {
          fields_rename_processor {
            description = "test"
            enabled     = true
            id          = "processor_test_logs_#name#"
            matcher     = "true"
            field {
              from_name = "#name#"
              to_name   = "#name#bar"
            }
          }
        }
        processor {
          technology_processor {
            enabled       = false
            id            = "processor_custom_technology_1_#name#"
            technology_id = "node_js_001"
            custom_matcher = "true"
            sample_data = "{}"
          }
        }
      }
    }
  }
}
