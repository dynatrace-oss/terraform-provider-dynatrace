resource "dynatrace_openpipeline_v2_bizevents_ingestsources" "ingest-source" {
  enabled = true
  display_name = "ingest-source"
  path_segment = "ingestsource.path.tf.#name#"
  static_routing {
    pipeline_type = "custom"
    pipeline_id = dynatrace_openpipeline_v2_bizevents_pipelines.pipeline.id
  }
  source_type = "http"
}

resource "dynatrace_openpipeline_v2_bizevents_pipelines" "pipeline" {
  display_name = "Pipeline"
  custom_id = "pipeline_1234_tf_#name#"
}
