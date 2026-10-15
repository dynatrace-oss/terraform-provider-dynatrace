# contains two configs
data "dynatrace_openpipeline_v2_bizevents_pipelines" "pipelines" {
  display_name = dynatrace_openpipeline_v2_bizevents_pipelines.pipeline1.display_name
  group_role   = "memberPipeline"

  depends_on = [dynatrace_openpipeline_v2_bizevents_pipelines.pipeline2, dynatrace_openpipeline_v2_bizevents_pipelines.pipeline3]
}

resource "dynatrace_openpipeline_v2_bizevents_pipelines" "pipeline1" {
  display_name = "#name#"
  custom_id    = "#name#.1"
  group_role   = "memberPipeline"
}

resource "dynatrace_openpipeline_v2_bizevents_pipelines" "pipeline2" {
  display_name = "#name#"
  custom_id    = "#name#.2"
  group_role   = "memberPipeline"
}

resource "dynatrace_openpipeline_v2_bizevents_pipelines" "pipeline3" {
  display_name = "#name#"
  custom_id    = "#name#.3"
  group_role   = "basePipeline"
  routing      = "notRoutable"
}
