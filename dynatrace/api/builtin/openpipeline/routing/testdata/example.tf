resource "dynatrace_openpipeline_v2_routing" "routing" {
  kind = "events"
  routing_entry {
    pipeline_type = "custom"
    pipeline_id = "pipelineId"
    matcher = "true"
    description = "test1"
  }
  routing_entry {
    enabled     = false
    pipeline_type = "custom"
    pipeline_id = "pipelineId"
    matcher = "false"
    description = "test2"
  }
  routing_entry {
    enabled     = false
    pipeline_type = "builtin"
    builtin_pipeline_id = "builtin_pipelineId"
    matcher = "true"
    description = "test3"
  }
}