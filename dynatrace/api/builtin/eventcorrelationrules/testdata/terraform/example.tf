resource "dynatrace_event_correlation_rules" "example" {
  display_name          = "example"
  enabled               = true
  matching_condition    = "event.kind == \"DAVIS_EVENT\" and event.category == \"AVAILABILITY\""
  grouping_fields       = ["dt.entity.host", "k8s.cluster.name"]
  time_window           = "5m"
  correlation_namespace = "namespace"
}
