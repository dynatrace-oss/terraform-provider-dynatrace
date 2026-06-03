variable "DT_GCP_TEST_UNIMPERSONABLE_SERVICE_ACCOUNT" {
    type        = string
    description = "The service account that should be used for the auth failure test. This service account must not allow impersonation by the DT GCP principal."
    sensitive = true
}

resource "dynatrace_gcp_connection" "auth_failure" {
  name = "#name#"
  type = "serviceAccountImpersonation"
  service_account_impersonation {
    service_account_id = var.DT_GCP_TEST_UNIMPERSONABLE_SERVICE_ACCOUNT
    consumers = [
      "SVC:com.dynatrace.da"
    ]
  }

  # Bound the create retry: the service account can never authenticate, so without a short timeout
  # the permanently-failing impersonation would be retried for the full DefaultCreateTimeout.
  timeouts {
    create = "10s"
  }
}
