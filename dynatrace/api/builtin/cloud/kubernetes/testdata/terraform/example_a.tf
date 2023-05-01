resource "dynatrace_kubernetes" "#name#" {
  enabled                            = true
  cloud_application_pipeline_enabled = true
  cluster_id                         = "#name#"
  cluster_id_enabled                 = true
  event_processing_active            = false
  label                              = "#name#"
  open_metrics_builtin_enabled       = false
  open_metrics_pipeline_enabled      = false
  pvc_monitoring_enabled             = false
  scope                              = "KUBERNETES_CLUSTER-1234567890000000"
}