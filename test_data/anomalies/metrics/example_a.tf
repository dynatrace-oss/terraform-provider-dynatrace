resource "dynatrace_metric_events" "#name#" {
  enabled                    = true
  event_entity_dimension_key = "dt.entity.azure_vm"
  summary                    = "Azure CPU"
  event_template {
    description = "The {metricname} value of {severity} was {alert_condition} your custom threshold of {threshold}."
    davis_merge = true
    event_type  = "RESOURCE"
    title       = "Azure CPU"
  }
  model_properties {
    type               = "STATIC_THRESHOLD"
    alert_condition    = "ABOVE"
    alert_on_no_data   = false
    dealerting_samples = 5
    samples            = 5
    threshold          = 40
    violating_samples  = 5
  }
  query_definition {
    type        = "METRIC_KEY"
    aggregation = "MAX"
    metric_key  = "builtin:cloud.azure.vm.cpuUsage"
    entity_filter {
      dimension_key = "dt.entity.azure_vm"
      conditions {
        condition {
          type     = "NAME"
          operator = "EQUALS"
          value    = "easytraveldemo-backend1"
        }
      }
    }
  }
}
