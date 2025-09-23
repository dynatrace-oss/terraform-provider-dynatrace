resource "dynatrace_openpipeline_v2_events_security_ingestsources" "minimal-source" {
  display_name = "min-ingest-source"
  enabled = true
  path_segment = "processor.ingestsource.path.tf.min.#name#"
  processing {}
}
