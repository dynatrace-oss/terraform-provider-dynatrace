resource "dynatrace_openpipeline_v2_events_security_routing" "routing" {
  routing_entries {
    routing_entry {
      enabled             = true
      pipeline_type       = "custom"
      pipeline_id = dynatrace_openpipeline_v2_events_security_pipelines.pipeline.id
      matcher             = "not matchesPhrase(record.title, \"Warning\")"
      description         = "Default route"
    }
  }
}

resource "dynatrace_openpipeline_v2_events_security_pipelines" "pipeline" {
  display_name = "Minimal pipeline"
  custom_id = "pipeline_Minimal_pipeline_1234_tf_#name#"
}
