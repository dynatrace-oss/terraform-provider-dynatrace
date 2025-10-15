resource "dynatrace_openpipeline_v2_security_events_pipelines" "example1" {
  display_name = "#name#"
  custom_id = "#name#"
  processing {
    processors {
      processor {
        type        = "drop"
        id          = "processor_Drop_unnecessary_records_3802"
        description = "Drop unnecessary records"
        matcher     = "not matchesPhrase(record.name, \"Warning\")"
        enabled     = true
      }
      processor {
        type        = "fieldsAdd"
        id          = "processor_Add_warning_flag_5434"
        description = "Add warning flag"
        matcher     = "matchesPhrase(record.name, \"Warning\")"
        sample_data = "{\n  \"record.name\": \"Warning record\" \n}"
        fields_add {
          fields {
            field {
              name  = "is_warning"
              value = "true"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "fieldsRemove"
        id          = "processor_Remove_details_field_8539"
        description = "Remove details field"
        sample_data = "{\n  \"record.name\": \"Warning\",\n  \"record.details\": \"some record details\"\n}"
        matcher     = "isNotNull(record.details)"
        fields_remove {
          fields = ["record.details"]
        }
        enabled     = true
      }
      processor {
        type        = "fieldsRename"
        id          = "processor_Rename_name_to_title_8530"
        description = "Rename name to title"
        sample_data = "{\n  \"record.name\": \"Warning\"\n}"
        matcher     = "true"
        fields_rename {
          fields {
            field {
              from_name = "record.name"
              to_name   = "record.title"
            }
          }
        }
        enabled     = true
      }
      processor {
        type        = "dql"
        id          = "processor_Combine_title_and_summary_to_name_8808"
        description = "Combine title and summary to name"
        sample_data = "{\n  \"record.title\": \"Warning\",\n  \"record.summary\": \"Request failed\"\n}"
        matcher     = "true"
        dql {
          script = "fieldsAdd record.name = concat(record.title, \" - \", record.summary)"
        }
        enabled     = true
      }
    }
  }
  davis {}
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


resource "dynatrace_openpipeline_v2_security_events_pipelines" "example2" {
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

resource "dynatrace_openpipeline_v2_security_events_pipelinegroups" "example" {
  display_name = "#name#"
  included_pipelines {
    included_pipeline {
      is_target_pipeline_placeholder = false
      pipeline_id = dynatrace_openpipeline_v2_security_events_pipelines.example1.id
      enabled_stages = ["processing", "metricExtraction"]
    }
  }
  targeted_pipelines = [dynatrace_openpipeline_v2_security_events_pipelines.example2.id]
}
