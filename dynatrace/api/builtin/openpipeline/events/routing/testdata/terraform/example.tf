resource "dynatrace_openpipeline_v2_events_routing" "routing" {
  routing_entries {
    routing_entry {
      enabled       = true
      pipeline_type = "custom"
      pipeline_id   = dynatrace_openpipeline_v2_events_pipelines.routing-pipeline.id
      matcher       = "matchesPhrase(record.title, \"Warning\")"
      description   = "Warning route"
    }
    routing_entry {
      enabled             = true
      pipeline_type       = "builtin"
      builtin_pipeline_id = "default"
      matcher             = "not matchesPhrase(record.title, \"Warning\")"
      description         = "Default route"
    }
  }
}

resource "dynatrace_openpipeline_v2_events_pipelines" "routing-pipeline" {
  display_name = "Routing pipeline"
  custom_id = "pipeline_Routing_pipeline_1234_tf_#name#"
  processing {}
  davis {}
  metric_extraction {}
  security_context {}
  cost_allocation {}
  product_allocation {}
  storage {}
  data_extraction {}
}