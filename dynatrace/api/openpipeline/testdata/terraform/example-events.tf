resource "dynatrace_openpipeline_events" "events" {
  endpoints {
    endpoint {
      enabled        = true
      default_bucket = "default_events"
      display_name   = "Custom ingest source"
      segment        = "something"
      routing {
        type        = "static"
        pipeline_id = "default"
      }
      processors {
        processor {
          fields_add_processor {
            description = "Custom add field"
            enabled     = true
            id          = "processor_Add_a_field_1_#name#"
            matcher     = "true"
            field {
              name  = "field"
              value = "value"
            }
          }
        }
        processor {
          fields_rename_processor {
            description = "Custom rename field"
            enabled     = true
            id          = "processor_Custom_rename_field_1_#name#"
            matcher     = "true"
            field {
              from_name = "new"
              to_name   = "old"
            }
          }
        }
        processor {
          fields_remove_processor {
            description = "Custom remove field"
            enabled     = true
            fields      = [ "field" ]
            id          = "processor_Custom_remove_field_1_#name#"
            matcher     = "true"
          }
        }
        processor {
          dql_processor {
            description = "Custom DQL"
            enabled     = true
            dql_script  = "fieldsAdd (\"test\")"
            id          = "processor_Custom_DQL_1_#name#"
            matcher     = "true"
          }
        }
        processor {
          drop_processor {
            description = "Custom drop processor"
            enabled     = true
            id          = "processor_custom_drop_1_#name#"
            matcher     = "true"
            sample_data = "{}"
          }
        }
      }
    }
  }
  pipelines {
    pipeline {
      enabled      = true
      display_name = "Custom pipeline 1"
      id           = "pipeline_Pipeline_1_#name#"
      processing {
        processor {
          fields_add_processor {
            description = "Add a field 1"
            enabled     = true
            id          = "processor_Add_a_field_2_#name#"
            matcher     = "true"
            field {
              name  = "field"
              value = "value"
            }
          }
        }
      }
    }
    pipeline {
      enabled      = true
      display_name = "Custom pipeline 2"
      id           = "pipeline_Pipeline_2_#name#"
      data_extraction {
        processor {
          davis_event_extraction_processor {
            description = "Custom event"
            enabled     = true
            id          = "processor_Custom_event_1_#name#"
            matcher     = "true"
            properties {
              key   = "event.type"
              value = "CUSTOM_ALERT"
            }
            properties {
              key   = "event.name"
              value = "test"
            }
          }
        }
      }
      metric_extraction {
        processor {
          value_metric_extraction_processor {
            description = "Custom value metric extraction"
            enabled     = true
            dimensions  = [ "availability" ]
            field       = "field1"
            id          = "processor_Custom_metric_extraction_1_#name#"
            matcher     = "true"
            metric_key  = "events.custom"
          }
        }
        processor {
          counter_metric_extraction_processor {
            description = "Custom counter metric extraction"
            enabled     = true
            id          = "processor_Custom_counter_metric_extraction_1_#name#"
            matcher     = "true"
            metric_key  = "events.counter"
          }
        }
      }

      processing {
        processor {
          fields_add_processor {
            description = "Custom add field"
            enabled     = true
            id          = "processor_Add_a_field_3_#name#"
            matcher     = "true"
            field {
              name  = "field"
              value = "value"
            }
          }
        }
        processor {
          fields_rename_processor {
            description = "Custom rename field"
            enabled     = true
            id          = "processor_Custom_rename_field_2_#name#"
            matcher     = "true"
            field {
              from_name = "new"
              to_name   = "old"
            }
          }
        }
        processor {
          fields_remove_processor {
            description = "Custom remove field"
            enabled     = true
            fields      = [ "field" ]
            id          = "processor_Custom_remove_field_2_#name#"
            matcher     = "true"
          }
        }
        processor {
          dql_processor {
            description = "Custom DQL"
            enabled     = true
            dql_script  = "fieldsAdd (\"test\")"
            id          = "processor_Custom_DQL_2_#name#"
            matcher     = "true"
          }
        }
        processor {
          drop_processor {
            description = "Custom drop processor"
            enabled     = true
            id          = "processor_custom_drop_2_#name#"
            matcher     = "true"
            sample_data = "{}"
          }
        }
      }
      security_context {
        processor {
          security_context_processor {
            description = "Custom security context 1"
            enabled     = true
            id          = "processor_Custom_security_context_1_#name#"
            matcher     = "true"
            sample_data = "{}"
            value {
              type     = "constant"
              constant = "string"
            }
          }
        }
        processor {
          security_context_processor {
            description = "Custom security context 2"
            enabled     = true
            id          = "processor_Custom_security_context_2_#name#"
            matcher     = "true"
            sample_data = "{}"
            value {
              type  = "field"
              field = "fieldname"
            }
          }
        }
        processor {
          security_context_processor {
            description = "Custom security context 3"
            enabled     = true
            id          = "processor_Custom_security_context_3_#name#"
            matcher     = "true"
            sample_data = "{}"
            value {
              type  = "multiValueConstant"
              multi_value_constant = ["multi", "value"]
            }
          }
        }
      }
      storage {
        catch_all_bucket_name = "default_events"
        processor {
          bucket_assignment_processor {
            description = "Custom bucket assignment"
            enabled     = true
            bucket_name = "default_events"
            id          = "processor_Custom_bucket_assignment_1_#name#"
            matcher     = "true"
            sample_data = "{}"
          }
        }
        processor {
          no_storage_processor {
            description = "Custom no storage assignment"
            enabled     = true
            id          = "processor_Custom_no_storage_assignment_1_#name#"
            matcher     = "true"
            sample_data = "{}"
          }
        }
      }
    }
  }
  routing {
    entry {
      enabled     = true
      matcher     = "true"
      note        = "Custom route"
      pipeline_id = "pipeline_Pipeline_1_#name#"
    }
  }
}
