# contains two configs
data "dynatrace_openpipeline_v2_bizevents_pipelines" "pipelines" {
  display_name = dynatrace_openpipeline_v2_bizevents_pipelines.pipeline1.display_name

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