resource "dynatrace_openpipeline_v2_smartscape_events_pipelines" "example1" {
  display_name = "#name#"
  custom_id    = "#name#"
  group_role   = "compositionPipeline"
  routing      = "notRoutable"
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
}


resource "dynatrace_openpipeline_v2_smartscape_events_pipelines" "example2" {
  display_name = "#name#-2"
  custom_id    = "#name#-2"
  group_role   = "memberPipeline"
}


resource "dynatrace_openpipeline_v2_smartscape_events_pipelinegroups" "example" {
  display_name = "#name#"
  composition {
    pipeline_group_composition {
      is_pipeline_placeholder = true
    }
    pipeline_group_composition {
      is_pipeline_placeholder = false
      stages {
        type    = "include"
        include = ["processing"]
      }
      pipeline_id = dynatrace_openpipeline_v2_smartscape_events_pipelines.example1.id
    }
  }
  member_stages {
    include = ["processing"]
    type    = "include"
  }
  member_pipelines = [dynatrace_openpipeline_v2_smartscape_events_pipelines.example2.id]
}
