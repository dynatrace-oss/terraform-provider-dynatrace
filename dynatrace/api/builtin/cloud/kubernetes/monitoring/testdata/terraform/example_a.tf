resource "dynatrace_k8s_monitoring" "#name#" {
  cloud_application_pipeline_enabled = true
  event_processing_active            = true
  filter_events                      = true
  include_all_fdi_events             = true
  open_metrics_builtin_enabled       = false
  open_metrics_pipeline_enabled      = true
  scope                              = "KUBERNETES_CLUSTER-1234567890000000"
  event_patterns {
    event_pattern {
      active  = true
      label   = "Node events"
      pattern = "involvedObject.kind=Node"
    }
  }
}