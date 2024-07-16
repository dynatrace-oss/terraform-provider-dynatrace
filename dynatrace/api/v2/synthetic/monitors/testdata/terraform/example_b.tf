resource "dynatrace_network_monitor" "ICMP_Test" {
  name          = "ICMP Test"
  type          = "MULTI_PROTOCOL"
  enabled       = false
  frequency_min = 30
  locations     = [ "SYNTHETIC_LOCATION-39F97465BE7BF644" ]
  outage_handling {
    global_consecutive_outage_count_threshold = 1
    global_outages                            = true
  }
  steps {
    step {
      name         = "ICMP Test"
      request_type = "ICMP"
      target_list  = [ "8.8.8.8" ]
      properties = {
        "ICMP_PACKET_SIZE"          = "32"
        "ICMP_NUMBER_OF_PACKETS"    = "1"
        "ICMP_TIMEOUT_FOR_REPLY"    = "PT1S"
        "ICMP_TIME_TO_LIVE"         = "64"
        "ICMP_DO_NOT_FRAGMENT_DATA" = "true"
        "ICMP_IP_VERSION"           = "4"
        "EXECUTION_TIMEOUT"         = "PT1S"
      }
      constraints {
        constraint {
          type = "SUCCESS_RATE_PERCENT"
          properties = {
            "value"    = "80"
            "operator" = ">="
          }
        }
      }
      request_configurations {
        request_configuration {
          constraints {
            constraint {
              type = "ICMP_SUCCESS_RATE_PERCENT"
              properties = {
                "value"    = "75"
                "operator" = ">="
              }
            }
          }
        }
      }
    }
  }
}