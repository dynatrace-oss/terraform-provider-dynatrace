data "dynatrace_openpipeline_v2_system_events_ingestsources" "ingestsources" {
  display_name = dynatrace_openpipeline_v2_system_events_ingestsources.source.display_name
  path_segment = "segment.#name#"
  source_type  = "http"
}

resource "dynatrace_openpipeline_v2_system_events_ingestsources" "source" {
  display_name = "#name#"
  enabled      = true
  path_segment = "segment.#name#"
  source_type  = "http"
}
