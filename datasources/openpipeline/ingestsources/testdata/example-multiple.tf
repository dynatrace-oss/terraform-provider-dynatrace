# should contain 2 items

data "dynatrace_openpipeline_v2_events_ingestsources" "ingestsources" {
  display_name = "#name#"
  source_type  = "http"
  depends_on   = [dynatrace_openpipeline_v2_events_ingestsources.source, dynatrace_openpipeline_v2_events_ingestsources.source2]
}

resource "dynatrace_openpipeline_v2_events_ingestsources" "source" {
  count        = 2
  display_name = "#name#"
  enabled      = true
  path_segment = "segment.#name#.${count.index}"
  source_type  = "http"
}

resource "dynatrace_openpipeline_v2_events_ingestsources" "source2" {
  display_name = "#name#-3"
  enabled      = true
  path_segment = "segment.#name#.3"
  source_type  = "http"
}
