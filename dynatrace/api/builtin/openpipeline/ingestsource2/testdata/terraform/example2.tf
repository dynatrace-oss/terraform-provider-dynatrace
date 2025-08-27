resource "dynatrace_openpipeline_v2_ingest_source2" "ingestsource" {
  kind = "events"
  default_bucket = "bucket"
  display_name = "ingest-source-with-processor"
  enabled = true
  path_segment = "processor.ingestsource.path"
  static_routing {
    pipeline_id = "pipelineId"
    pipeline_type = "custom"
  }
  processing {
    processor {
      drop_processor {
        enabled     = true
        id          = "proc-1"
        description = "my-proc-1"
        sample_data = "my sample data"
      }
    }
    processor {
      dql_processor {
        enabled     = true
        id          = "proc-2"
        description = "my-proc-2"
        sample_data = "my sample data"
        dql {
          script = "fieldsAdd true"
        }
      }
    }
    processor {
      fields_add_processor {
        enabled     = true
        id          = "proc-3"
        description = "my-proc-3"
        sample_data = "my sample data"
        fields_add {
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
    processor {
      fields_rename_processor {
        enabled     = true
        id          = "proc-4"
        description = "my-proc-4"
        sample_data = "my sample data"
        fields_rename {
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

    processor {
      fields_remove_processor {
        enabled     = true
        id          = "proc-5"
        description = "my-proc-5"
        sample_data = "my sample data"
        fields_remove {
          fields = ["to-remove-1", "to-remove-2"]
        }
      }
    }
  }
}
