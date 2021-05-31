resource "dynatrace_host_anomalies" "#name#" {
  connections {
    enabled                       = true
    enabled_on_graceful_shutdowns = false
  }
  cpu {
    enabled = true
  }
  disks {
    inodes {
      enabled = true
    }
    space {
      enabled = true
    }
    speed {
      enabled = true
    }
  }
  gc {
    enabled = true
  }
  java {
    out_of_memory {
      enabled = true
    }
    out_of_threads {
      enabled = true
    }
  }
  memory {
    enabled = true
    thresholds {
      linux {
        page_faults = 20
        usage       = 80
      }
      windows {
        page_faults = 100
        usage       = 80
      }
    }
  }
  network {
    connectivity {
      enabled = true
    }
    dropped_packets {
      enabled = true
    }
    errors {
      enabled = true
    }
    retransmission {
      enabled = true
    }
    utilization {
      enabled = true
    }
  }
}
