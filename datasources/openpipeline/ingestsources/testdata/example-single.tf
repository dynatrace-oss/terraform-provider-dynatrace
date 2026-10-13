# should contain one item

data "dynatrace_openpipeline_v2_events_ingestsources" "ingestsources" {
  path_segment = "segment.#name#"
  source_type  = "http"
  display_name = "#name#"
  depends_on   = [dynatrace_openpipeline_v2_events_ingestsources.source, dynatrace_openpipeline_v2_events_ingestsources.source2]
}

resource "dynatrace_openpipeline_v2_events_ingestsources" "source" {
  display_name = "#name#"
  enabled      = true
  path_segment = "segment.#name#"
  source_type  = "http"
}

resource "dynatrace_openpipeline_v2_events_ingestsources" "source2" {
  display_name = "#name#"
  enabled      = true
  path_segment = "segment.#name#.2"
  source_type  = "http"
}
