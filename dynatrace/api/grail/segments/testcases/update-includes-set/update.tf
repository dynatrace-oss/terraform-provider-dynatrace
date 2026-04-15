resource "dynatrace_segment" "segment" {
  name        = "#name#"
  description = "Example description"
  is_public   = true
  includes {
    items {
      data_object = "dt.entity.cloud_application"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    # updated
    items {
      data_object = "dt.entity.service"
      filter      = ""
      relationship {
        name   = "runs_on"
        target = "dt.entity.host"
      }
    }
  }
}
