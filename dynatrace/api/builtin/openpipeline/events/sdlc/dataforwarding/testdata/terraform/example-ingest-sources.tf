resource "dynatrace_openpipeline_v2_events_sdlc_dataforwarding" "example" {
  forwarding_name   = "#name#"
  enabled           = false
  matcher           = "true"
  cloud_vendor_type = "gcp"
  gcp_connection {
    bucket_name   = "my-bucket"
    connection_id = dynatrace_gcp_connection.connection.id
  }
  data_forwarding_type = "raw"
  ingest_sources       = [dynatrace_openpipeline_v2_events_sdlc_ingestsources.source.id]
  bulk_pattern         = "<YYYYMMDD>/<HH>/<HHmmss.SSSS>_<bulk-id>.json.gz"
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
        enabled = true
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
        enabled = true
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
        enabled = true
      }
    }
  }
}

resource "dynatrace_openpipeline_v2_events_sdlc_ingestsources" "source" {
  display_name = "min-ingest-source"
  enabled      = true
  path_segment = "processor.ingestsource.path.tf.min.#name#"
}


# Create GCP connection
variable "DT_GCP_TEST_IMPERSONABLE_SERVICE_ACCOUNT" {
  type        = string
  description = "The service account that should be used for the GCP connection setup end-to-end test."
  sensitive   = true
}

resource "dynatrace_gcp_connection" "connection" {
  name = "#name#"
  type = "serviceAccountImpersonation"
  service_account_impersonation {
    service_account_id = var.DT_GCP_TEST_IMPERSONABLE_SERVICE_ACCOUNT
    consumers = [
      "SVC:com.dynatrace.openpipeline"
    ]
  }
}
