resource "dynatrace_openpipeline_v2_events_security_pipelines" "min-pipeline" {
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
