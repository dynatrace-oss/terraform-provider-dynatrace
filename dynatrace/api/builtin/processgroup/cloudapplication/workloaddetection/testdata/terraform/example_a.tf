# ID vu9U3hXa3q0AAAABADpidWlsdGluOnByb2Nlc3MtZ3JvdXAuY2xvdWQtYXBwbGljYXRpb24td29ya2xvYWQtZGV0ZWN0aW9uAAZ0ZW5hbnQABnRlbmFudAAkYjcwNmY4NWYtNWFkNC0zY2ZmLWJhYzMtZDg4YzFmNTkzMjgwvu9U3hXa3q0
resource "dynatrace_cloudapp_workloaddetection" "cloud_app_workload_detection" {
  cloud_foundry {
    enabled = false
  }
  docker {
    enabled = true
  }
  kubernetes {
    enabled = true
    filters {
      filter {
        enabled = false
        inclusion_toggles {
          inc_basepod   = false
          inc_container = true
          inc_namespace = true
          inc_product   = true
          inc_stage     = true
        }
        match_filter {
          match_operator = "EXISTS"
        }
      }
    }
  }
}
