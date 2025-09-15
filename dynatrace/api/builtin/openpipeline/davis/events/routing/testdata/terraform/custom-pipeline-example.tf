resource "dynatrace_openpipeline_v2_davis_events_routing" "routing" {
  routing_entries {
    routing_entry {
      enabled             = true
      pipeline_type       = "custom"
      pipeline_id = dynatrace_openpipeline_v2_davis_events_pipelines.pipeline.id
      matcher             = "not matchesPhrase(record.title, \"Warning\")"
      description         = "Default route"
    }
  }
}

resource "dynatrace_openpipeline_v2_davis_events_pipelines" "pipeline" {
  display_name = "Minimal pipeline"
  custom_id = "pipeline_Minimal_pipeline_1234_tf_#name#"
  processing {}
  davis {}
  metric_extraction {}
  security_context {}
  cost_allocation {}
  product_allocation {}
  storage {}
  data_extraction {}
}
