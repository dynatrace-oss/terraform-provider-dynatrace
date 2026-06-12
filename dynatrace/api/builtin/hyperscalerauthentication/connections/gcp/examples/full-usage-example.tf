variable "gcp_project_id" {
  type        = string
  description = "The Google Cloud Platform project ID"
}

# Configure the GCP provider
provider "google" {
  project = var.gcp_project_id
}

# Create a Service Account in GCP
resource "google_service_account" "impersonable_service_account" {
  account_id   = "test-service-account"
  display_name = "Test Service Account"
}

# Provision the DT GCP Principal. This resource takes no input; it triggers creation of the
# Dynatrace-managed principal (a singleton) and exposes it via the `principal` attribute.
resource "dynatrace_gcp_principal" "principal" {}

# Grant DT GCP Principal access to the service account
resource "google_service_account_iam_member" "wif_binding" {
  service_account_id = google_service_account.impersonable_service_account.name
  member             = "serviceAccount:${dynatrace_gcp_principal.principal.principal}"
  role               = "roles/iam.serviceAccountTokenCreator"
}

# Create GCP connection
resource "dynatrace_gcp_connection" "my_gcp_connection" {
  name = "My Gcp Connection"
  type = "serviceAccountImpersonation"
  service_account_impersonation {
    service_account_id = google_service_account.impersonable_service_account.email
    consumers = [
      "SVC:com.dynatrace.da"
    ]
  }

  depends_on = [
    google_service_account_iam_member.wif_binding
  ]

  # The create timeout is set to 2 minutes by default.
  # If you wish to adjust this, you can do so using a `timeouts` block:
  # timeouts {
  #   create = "3m"
  # }
}
