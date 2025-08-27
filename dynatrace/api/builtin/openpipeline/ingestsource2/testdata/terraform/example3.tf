resource "dynatrace_openpipeline_v3_azuge_log_forwarding_ingest_source" "ingestsource" {
  default_bucket = "bucket"
  display_name = "ingest-source-with-processor"
  enabled = true
  path_segment = "processor.ingestsource.path"
  static_routing {
    pipeline_id = "pipelineId"
    pipeline_type = "custom"
  }
  processing {
    processors {
      type        = "drop"
      enabled     = true
      id          = "proc-1"
      description = "my-proc-1"
      sample_data = "my sample data"
      matcher     = "not true"
    }
    processors {
      type        = "dql"
      enabled     = true
      id          = "proc-2"
      description = "my-proc-2"
      sample_data = "my sample data"
      matcher     = "not true"
      processor {
        description = ""
        enabled     = false
        id          = ""
        type        = "dql"
        dql {
          script = "fieldsAdd true"
        }
      }
    }
    processors {
      type        = "fieldsAdd"
      enabled     = true
      id          = "proc-3"
      description = "my-proc-3"
      sample_data = "my sample data"
      matcher     = "not true"
      processor {
        description = ""
        enabled     = false
        id          = ""
        type        = "fieldsAdd"
        fields_add {
          fields {
            field {
              name = "some-name"
              value = "some-value"
            }
            field {
              name = "some-other-name"
              value = "some-other-value"
            }
          }
        }
      }
    }
    processors {
      type        = "fieldsRename"
      enabled     = true
      id          = "proc-4"
      description = "my-proc-4"
      sample_data = "my sample data"
      matcher     = "not true"
      processor {
        description = ""
        enabled     = false
        id          = ""
        type        = "fieldsRename"
        fields_rename {
          fields {
            field {
              from_name = "from-name"
              to_name = "to-name"
            }
            field {
              from_name = "from-other-name"
              to_name = "to-other-name"
            }
          }
        }
      }
    }

    processors {
      type        = "fieldsRemove"
      enabled     = true
      id          = "proc-5"
      description = "my-proc-5"
      sample_data = "my sample data"
      matcher     = "not true"
      processor {
        description = ""
        enabled     = false
        id          = ""
        type        = "fieldsRemove"
        fields_remove {
          fields = ["to-remove-1", "to-remove-2"]
        }
      }
    }
  }
}
