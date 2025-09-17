resource "dynatrace_openpipeline_v2_metrics_ingestsources" "maximal-source" {
  enabled = true
  display_name = "max-ingestsource"
  path_segment = "processor.ingestsource.path.max.tf.#name#"
  static_routing {
    pipeline_type = "builtin"
    builtin_pipeline_id = "default"
  }
  default_bucket = "default_events"
  processing {
    processors {
      processor {
        enabled     = true
        type        = "drop"
        id          = "processor_Drop_unnecessary_records_1234"
        description = "Drop unnecessary records"
        matcher     = "not matchesPhrase(record.name, \"Error\") and not matchesPhrase(record.name, \"Warning\")"
      }
      processor {
        enabled     = true
        type = "fieldsAdd"
        id          = "processor_Add_error_flag_6132"
        description = "Add error flag"
        matcher = "matchesPhrase(record.name, \"Error\")"
        sample_data = "{\n  \"record.name\": \"Error record\" \n}"
        fields_add {
          fields {
            field {
              name = "is_error"
              value = "true"
            }
          }
        }
      }
      processor {
        enabled     = true
        type = "fieldsRemove"
        id          = "processor_Remove_details_field_8919"
        description = "Remove details field"
        matcher = "isNotNull(record.details)"
        sample_data = "{\n  \"record.name\": \"Error\",\n  \"record.details\": \"some record details\"\n}"
        fields_remove {
          fields = ["record.details"]
        }
      }
      processor {
        enabled     = true
        type = "fieldsRename"
        id          = "processor_Rename_name_to_title_5347"
        description = "Rename name to title"
        matcher = "true"
        sample_data = "{\n  \"record.name\": \"Error\"\n}"
        fields_rename {
          fields {
            field {
              from_name = "record.name"
              to_name = "record.title"
            }
          }
        }
      }
      processor {
        enabled     = true
        type        = "dql"
        id          = "processor_Combine_title_and_summary_to_name_1244"
        description = "Combine title and summary to name"
        sample_data = "{\n  \"record.title\": \"Error\",\n  \"record.summary\": \"Request failed\"\n}"
        matcher     = "true"
        dql {
          script = "fieldsAdd record.name = concat(record.title, \" - \", record.summary)"
        }
      }
    }
  }
}
