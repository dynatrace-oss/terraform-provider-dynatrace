resource "dynatrace_event_correlation_rules" "first" {
  display_name          = "first"
  enabled               = true
  matching_condition    = "event.kind == \"DAVIS_EVENT\" and event.category == \"AVAILABILITY\""
  grouping_fields       = ["dt.entity.host", "k8s.cluster.name"]
  time_window           = "5m"
  correlation_namespace = "namespace"
}


resource "dynatrace_event_correlation_rules" "second" {
  display_name          = "second"
  enabled               = true
  matching_condition    = "event.kind == \"DAVIS_EVENT\" and event.category == \"AVAILABILITY\""
  grouping_fields       = ["dt.entity.host", "k8s.cluster.name"]
  time_window           = "5m"
  correlation_namespace = "namespace"
  insert_after          = dynatrace_event_correlation_rules.first.id
}
