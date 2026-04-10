resource "dynatrace_metric_events" "events" {
  enabled                    = true
  event_entity_dimension_key = "dt.entity.host"
  summary                    = "sample_name"
  event_template {
    description   = "some description"
    davis_merge = false
    event_type    = "CUSTOM_ALERT"
    title         = "sample_name"
    # Set
    metadata {
      metadata_key   = "key"
      metadata_value = "value"
    }
    metadata {
      metadata_key   = "keyEdit"
      metadata_value = "valuEdit"
    }
    metadata {
      metadata_key   = "keyNew"
      metadata_value = "valuNew"
    }
  }
  model_properties {
    type               = "STATIC_THRESHOLD"
    alert_condition    = "ABOVE"
    alert_on_no_data   = false
    dealerting_samples = 5
    samples            = 5
    threshold          = 1
    violating_samples  = 3
  }
  query_definition {
    type        = "METRIC_KEY"
    aggregation = "AVG"
    metric_key  = "builtin:host.disk.usedPct"
    # Set
    dimension_filter {
      filter {
        dimension_key   = "dt.entity.host"
        dimension_value = "HOST-0000000000000000"
      }
      filter {
        dimension_key   = "dt.entity.disk"
        dimension_value = "DISK-0000000000000001"
      }
    }
    entity_filter {
      dimension_key = "dt.entity.host"
      # Set
      conditions {
        condition {
          type     = "NAME"
          operator = "EQUALS"
          value    = "HOST-0000000000000000"
        }
        condition {
          type     = "ENTITY_ID"
          operator = "EQUALS"
          value    = "HOST-0000000000000001"
        }
        condition {
          type     = "HOST_GROUP_NAME"
          operator = "EQUALS"
          value    = "HOST-0000000000000000"
        }
      }
    }
  }
}
