resource "dynatrace_openpipeline_v2_security_events_routing" "routing" {
  routing_entries {
    routing_entry {
      enabled             = true
      pipeline_type       = "custom"
      pipeline_id = dynatrace_openpipeline_v2_security_events_pipelines.pipeline.id
      matcher             = "not matchesPhrase(record.title, \"Warning\")"
      description         = "Default route"
    }
  }
}

resource "dynatrace_openpipeline_v2_security_events_pipelines" "pipeline" {
  display_name = "Minimal pipeline"
  custom_id = "pipeline_Minimal_pipeline_1234_tf_#name#"
}
