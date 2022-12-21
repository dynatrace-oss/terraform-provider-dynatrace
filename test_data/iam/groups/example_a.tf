resource "dynatrace_iam_group" "#name#" {
  name          = "#name#"
  # description = ""
  permissions {
    permission {
      name  = "tenant-view-sensitive-request-data"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-configure-request-capture-data"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-logviewer"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-agent-install"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-manage-settings"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-viewer"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-replay-sessions-with-masking"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-replay-sessions-without-masking"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-manage-security-problems"
      type  = "tenant"
      scope = "siz65484"
    }
    permission {
      name  = "tenant-viewer"
      type  = "tenant"
      scope = "gix28347"
    }
    permission {
      name  = "tenant-manage-support-tickets"
      type  = "tenant"
      scope = "siz65484"
    }
  }
}
