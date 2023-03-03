resource "dynatrace_host_anomalies_v2" "#name#" {
  scope = "HOST-1234567890000000"
  host {
    connection_lost_detection {
      enabled               = true
      on_graceful_shutdowns = "DONT_ALERT_ON_GRACEFUL_SHUTDOWN"
    }
    high_cpu_saturation_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        cpu_saturation = 95
        event_thresholds {
          dealerting_evaluation_window = 30
          dealerting_samples           = 30
          violating_evaluation_window  = 30
          violating_samples            = 18
        }
      }
    }
    high_gc_activity_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        gc_suspension_percentage = 25
        gc_time_percentage       = 40
        event_thresholds {
          dealerting_evaluation_window = 30
          dealerting_samples           = 30
          violating_evaluation_window  = 30
          violating_samples            = 18
        }
      }
    }
    high_memory_detection {
      enabled        = true
      detection_mode = "auto"
    }
    high_system_load_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        system_load = 1
        event_thresholds {
          dealerting_evaluation_window = 30
          dealerting_samples           = 30
          violating_evaluation_window  = 6
          violating_samples            = 1
        }
      }
    }
    out_of_memory_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        out_of_memory_exceptions_number = 1
        event_thresholds {
          dealerting_evaluation_window = 30
          dealerting_samples           = 30
          violating_evaluation_window  = 6
          violating_samples            = 1
        }
      }
    }
    out_of_threads_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        out_of_threads_exceptions_number = 1
        event_thresholds {
          dealerting_evaluation_window = 30
          dealerting_samples           = 30
          violating_evaluation_window  = 6
          violating_samples            = 1
        }
      }
    }
  }
  network {
    high_network_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        errors_percentage = 90
        event_thresholds {
          dealerting_evaluation_window = 30
          dealerting_samples           = 30
          violating_evaluation_window  = 30
          violating_samples            = 18
        }
      }
    }
    network_dropped_packets_detection {
      enabled        = true
      detection_mode = "auto"
    }
    network_errors_detection {
      enabled        = true
      detection_mode = "auto"
    }
    network_high_retransmission_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        retransmission_rate_percentage          = 10
        retransmitted_packets_number_per_minute = 10
        event_thresholds {
          dealerting_evaluation_window = 30
          dealerting_samples           = 30
          violating_evaluation_window  = 30
          violating_samples            = 18
        }
      }
    }
    network_tcp_problems_detection {
      enabled        = true
      detection_mode = "auto"
    }
  }
}
