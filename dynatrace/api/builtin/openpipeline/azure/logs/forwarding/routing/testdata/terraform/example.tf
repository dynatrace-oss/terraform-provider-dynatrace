resource "dynatrace_openpipeline_v2_azure_logs_forwarding_routing" "routing" {
  routing_entries {
    routing_entry {
      enabled             = true
      pipeline_type       = "builtin"
      builtin_pipeline_id = "default"
      matcher             = "not matchesPhrase(record.title, \"Warning\")"
      description         = "Default route"
    }
  }
}