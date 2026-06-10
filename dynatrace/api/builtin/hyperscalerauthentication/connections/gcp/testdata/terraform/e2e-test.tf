variable "DT_GCP_TEST_IMPERSONABLE_SERVICE_ACCOUNT" {
  type        = string
  description = "The service account that should be used for the GCP connection setup end-to-end test."
  sensitive   = true
}

# Create GCP connection
resource "dynatrace_gcp_connection" "example" {
  name = "#name#"
  type = "serviceAccountImpersonation"
  service_account_impersonation {
    service_account_id = var.DT_GCP_TEST_IMPERSONABLE_SERVICE_ACCOUNT
    consumers = [
      "SVC:com.dynatrace.da"
    ]
  }
}