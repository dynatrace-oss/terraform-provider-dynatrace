resource "dynatrace_openpipeline_v2_metrics_ingestsources" "ingest-source" {
  enabled = true
  display_name = "ingest-source"
  path_segment = "ingestsource.path.tf.#name#"
  static_routing {
    pipeline_type = "custom"
    pipeline_id = dynatrace_openpipeline_v2_metrics_pipelines.pipeline.id
  }
  processing {
  }
}

resource "dynatrace_openpipeline_v2_metrics_pipelines" "pipeline" {
  display_name = "Pipeline"
  custom_id = "pipeline_1234_tf_#name#"
  processing {}
  davis {}
  metric_extraction {}
  security_context {}
  cost_allocation {}
  product_allocation {}
  storage {}
  data_extraction {}
}
