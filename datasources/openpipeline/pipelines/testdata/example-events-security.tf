data "dynatrace_openpipeline_v2_events_security_pipelines" "pipelines" {
  display_name = dynatrace_openpipeline_v2_events_security_pipelines.pipeline.display_name
  group_role   = "memberPipeline"
}

resource "dynatrace_openpipeline_v2_events_security_pipelines" "pipeline" {
  display_name = "#name#"
  custom_id    = "pipeline_Minimal_pipeline_1234_tf_#name#"
  group_role   = "memberPipeline"
}
