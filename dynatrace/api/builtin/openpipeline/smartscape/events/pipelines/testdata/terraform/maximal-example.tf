resource "dynatrace_openpipeline_v2_smartscape_events_pipelines" "max-pipeline" {
  display_name = "Warning pipeline"
  custom_id    = "pipeline_Warning_pipeline_2773_tf_#name#"
  metadata_list {
    metadata {
      entry_key   = "environment"
      entry_value = "production"
    }
  }
  processing {
    processors {
      processor {
        type        = "drop"
        id          = "processor_Drop_unnecessary_records_3802"
        description = "Drop unnecessary records"
        matcher     = "not matchesPhrase(name, \"Warning\")"
        enabled     = true
      }
      processor {
        type        = "fieldsAdd"
        id          = "processor_Add_warning_flag_5434"
        description = "Add warning flag"
        matcher     = "matchesPhrase(name, \"Warning\")"
        sample_data = "{\n  \"name\": \"Warning record\" \n}"
        fields_add {
          fields {
            field {
              name       = "custom.is_warning"
              value      = "true"
              value_type = "boolean"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "fieldsRemove"
        id          = "processor_Remove_details_field_8539"
        description = "Remove details field"
        sample_data = "{\n  \"name\": \"Warning\",\n  \"custom.details\": \"some record details\"\n}"
        matcher     = "isNotNull(custom.details)"
        fields_remove {
          fields = ["custom.details"]
        }
        enabled = true
      }
      processor {
        type        = "fieldsRename"
        id          = "processor_Rename_title_to_summary_8530"
        description = "Rename title to summary"
        sample_data = "{\n  \"custom.title\": \"Warning\"\n}"
        matcher     = "true"
        fields_rename {
          fields {
            field {
              from_name = "custom.title"
              to_name   = "custom.summary"
            }
          }
        }
        enabled = true
      }
      processor {
        type        = "dql"
        id          = "processor_Combine_summary_and_source_8808"
        description = "Combine summary and source"
        sample_data = "{\n  \"custom.summary\": \"Warning\",\n  \"custom.source\": \"Request failed\"\n}"
        matcher     = "true"
        dql {
          script = "fieldsAdd custom.combined = concat(custom.summary, \" - \", custom.source)"
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
            type                 = "multiValueConstant"
            multi_value_constant = ["ACME1", "ACME2"]
          }
        }
        enabled = true
      }
    }
  }
}
