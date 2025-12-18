resource "dynatrace_synthetic_location" "Test" {
  name                                  = "Test"
  auto_update_chromium                  = true
  availability_location_outage          = true
  availability_node_outage              = true
  availability_notifications_enabled    = true
  city                                  = "San Francisco de Asis"
  country_code                          = "VE"
  deployment_type                       = "STANDARD"
  latitude                              = 10.0756
  location_node_outage_delay_in_minutes = 3
  longitude                             = -67.5442
  nodes                                 = [ ]
  region_code                           = "04"
}

resource "time_sleep" "wait_for_location" {
  depends_on = [dynatrace_synthetic_location.Test]
  create_duration = "15s"
}

resource "dynatrace_network_monitor" "DNS_Test" {
  depends_on = [time_sleep.wait_for_location]
  name          = "DNS Test"
  description   = "This is an example DNS test"
  type          = "MULTI_PROTOCOL"
  enabled       = false
  frequency_min = 15
  locations     = [ dynatrace_synthetic_location.Test.id ]
  outage_handling {
    global_consecutive_outage_count_threshold = 1
    global_outages                            = true
  }
  performance_thresholds {
    enabled = true
    thresholds {
      threshold {
        aggregation        = "AVG"
        dealerting_samples = 5
        samples            = 5
        step_index         = 0
        threshold          = 100
        violating_samples  = 3
      }
    }
  }
  steps {
    step {
      name         = "DNS Test"
      request_type = "DNS"
      target_list  = [ "google.com", "yahoo.com" ]
      properties = {
        "DNS_RECORD_TYPES"  = "A"
        "EXECUTION_TIMEOUT" = "PT2S"
      }
      constraints {
        constraint {
          type = "SUCCESS_RATE_PERCENT"
          properties = {
            "value"    = "90"
            "operator" = ">="
          }
        }
        constraint {
          type = "SUCCESS_RATE_PERCENT"
          properties = {
            "value"    = "90"
            "operator" = "<"
          }
        }
      }
      request_configurations {
        request_configuration {
          constraints {
            constraint {
              type = "DNS_STATUS_CODE"
              properties = {
                "operator"   = "="
                "statusCode" = "0"
              }
            }
          }
        }
      }
    }
  }
  tags {
    tag {
      context = "CONTEXTLESS"
      key     = "Key1"
      source  = "USER"
      value   = "Value1"
    }
  }
}
