resource "dynatrace_log_events" "#name#" {
  enabled = false
  query   = "content=\"terraform\""
  summary = "Created by Terraform"
  event_template {
    description = "Created by Terraform"
    event_type  = "INFO"
    title       = "{content}"
    metadata {
      item {
        metadata_key   = "terraform.key"
        metadata_value = "terraform.value"
      }
    }
  }
}