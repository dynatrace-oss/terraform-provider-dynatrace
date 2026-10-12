data "dynatrace_openpipeline_v2_security_events_pipelines" "pipelines" {
  display_name = dynatrace_openpipeline_v2_security_events_pipelines.pipeline.display_name
}

resource "dynatrace_openpipeline_v2_security_events_pipelines" "pipeline" {
  display_name = "#name#"
  custom_id    = "pipeline_Minimal_pipeline_1234_tf_#name#"
}
