# contains two configs
data "dynatrace_openpipeline_v2_bizevents_pipelines" "pipelines" {
  custom_id = dynatrace_openpipeline_v2_bizevents_pipelines.pipeline1.custom_id

  depends_on = [dynatrace_openpipeline_v2_bizevents_pipelines.pipeline2]
}

resource "dynatrace_openpipeline_v2_bizevents_pipelines" "pipeline1" {
  display_name = "#name#"
  custom_id    = "#name#.1"
}

resource "dynatrace_openpipeline_v2_bizevents_pipelines" "pipeline2" {
  display_name = "#name#"
  custom_id    = "#name#.2"
}