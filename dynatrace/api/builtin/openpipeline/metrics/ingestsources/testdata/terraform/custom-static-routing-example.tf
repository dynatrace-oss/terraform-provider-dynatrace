resource "dynatrace_openpipeline_v2_metrics_ingestsources" "ingest-source" {
  enabled = true
  display_name = "ingest-source"
  path_segment = "ingestsource.path.tf.#name#"
  source_type = "http"
  static_routing {
    pipeline_type = "custom"
    pipeline_id = dynatrace_openpipeline_v2_metrics_pipelines.pipeline.id
  }
}

resource "dynatrace_openpipeline_v2_metrics_pipelines" "pipeline" {
  display_name = "Pipeline"
  custom_id = "pipeline_1234_tf_#name#"
}
