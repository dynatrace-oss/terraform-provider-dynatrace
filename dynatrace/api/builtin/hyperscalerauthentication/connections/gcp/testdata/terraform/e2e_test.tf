variable "gcp_project_number" {
  type        = string
  description = "The Google Cloud Platform project number"
}

terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.33.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.14.0"
    }
  }
}

# Configure the GCP provider
provider "google" {
  project = var.gcp_project_number
}

# Create a Service Account in GCP
resource "google_service_account" "gcp_service_account" {
  account_id   = "#name#"
  display_name = "#name#"
}

# Fetch the DT GCP Principal
data "dynatrace_gcp_principal" "principal" {
}

# Grant DT GCP Principal access to the service account
resource "google_service_account_iam_member" "wif_binding" {
  service_account_id = google_service_account.gcp_service_account.name
  member  = data.dynatrace_gcp_principal.principal.principal
  role    = "roles/iam.serviceAccountTokenCreator"
}

# Wait for IAM binding to propagate
resource time_sleep "wait_for_iam_binding" {
  depends_on = [google_service_account_iam_member.wif_binding]
  create_duration = "2m"
}

# Create GCP connection
resource "dynatrace_gcp_connection" "example" {
  name = "#name#"
  type = "serviceAccountImpersonation"
  service_account_impersonation {
    service_account_id  = google_service_account.gcp_service_account.name
    consumers = [
      "DA"
    ]
  }

  depends_on = [
    time_sleep.wait_for_iam_binding
  ]
}
