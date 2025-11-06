resource "dynatrace_openpipeline_v2_security_events_ingestsources" "minimal-source" {
  display_name = "min-ingest-source"
  enabled = true
  path_segment = "processor.ingestsource.path.tf.min.#name#"
}
