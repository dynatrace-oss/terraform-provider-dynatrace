data "dynatrace_openpipeline_v2_bizevents_pipelines" "pipelines" {
  custom_id = dynatrace_openpipeline_v2_bizevents_pipelines.pipeline.custom_id
}

resource "dynatrace_openpipeline_v2_bizevents_pipelines" "pipeline" {
  display_name = "#name#"
  custom_id    = "pipeline_Minimal_pipeline_1234_tf_#name#"
}