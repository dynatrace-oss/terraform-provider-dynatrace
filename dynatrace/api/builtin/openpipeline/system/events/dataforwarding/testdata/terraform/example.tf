resource "dynatrace_openpipeline_v2_system_events_dataforwarding" "example" {
  forwarding_name   = "#name#"
  enabled           = false
  matcher           = "true"
  cloud_vendor_type = "gcp"
  gcp_connection {
    bucket_name   = "my-bucket"
    connection_id = dynatrace_gcp_connection.connection.id
  }
  data_forwarding_type = "processed"
  pipelines            = [dynatrace_openpipeline_v2_system_events_pipelines.pipeline.id]
  bulk_pattern         = "<YYYYMMDD>/<HH>/<HHmmss.SSSS>_<bulk-id>.json.gz"
}

resource "dynatrace_openpipeline_v2_system_events_pipelines" "pipeline" {
  display_name = "Minimal pipeline"
  custom_id    = "pipeline_Minimal_pipeline_1234_tf_#name#"
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
