resource "dynatrace_log_events" "events" {
  enabled = false
  query   = "matchesPhrase(content, \"terratest\")"
  summary = "Created by Terraform"
  event_template {
    description = "Created by Terraform"
    event_type  = "INFO"
    title       = "{content}"
    metadata {
      item {
        metadata_key   = "terraform.key1"
        metadata_value = "terraform.value1"
      }
      item {
        metadata_key   = "terraform.keyEdit"
        metadata_value = "terraform.valueEdit"
      }
    }
  }
}
