resource "dynatrace_document" "sample_dashboard" {
  name = "#name#-dashboard"
  type = "dashboard"
  content = jsonencode({
    "annotations" : [],
    "importedWithCode" : false,
    "layouts" : {
      "0" : {
        "h" : 6,
        "w" : 8,
        "x" : 0,
        "y" : 0
      }
    },
    "settings" : {},
    "tiles" : {
      "0" : {
        "content" : "Hello World!",
        "type" : "markdown"
      }
    },
    "variables" : [],
    "version" : 21
  })
  private = true
}

resource "dynatrace_iam_service_user" "sample_service_user" {
  name        = "#name#-service-user"
  description = "Service user that can access the dashboard"
}

resource "dynatrace_iam_group" "sample_group1" {
  name        = "#name#-group1"
  description = "First group that can access the dashboard"
}

resource "dynatrace_direct_shares" "direct_share" {
  access      = "read-write"
  document_id = dynatrace_document.sample_dashboard.id
  recipients {
    # to update
    recipient {
      type = "user"
      id   = dynatrace_iam_service_user.sample_service_user.id
    }
    recipient {
      type = "group"
      id   = dynatrace_iam_group.sample_group1.id
    }
  }
}
