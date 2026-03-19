resource "dynatrace_direct_shares" "this" {
  document_id = dynatrace_document.this.id
  access      = "read-write"

  recipients {
    recipient {
      id   = dynatrace_iam_service_user.sample_service_user.id
      type = "user"
    }

    recipient {
      id   = dynatrace_iam_group.sample_group.id
      type = "group"
    }
  }
}

resource "dynatrace_document" "this" {
  type = "dashboard"
  name = "#name#"
  content = jsonencode(
    {
      "version" : 1,
      "variables" : [],
      "tiles" : {
        "0" : {
          "type" : "markdown",
          "title" : "",
          "content" : "Dashboard content"
        }
      },
      "layouts" : {
        "0" : {
          "x" : 0,
          "y" : 0,
          "w" : 24,
          "h" : 14
        }
      }
    }
  )
}

resource "dynatrace_iam_service_user" "sample_service_user" {
  name        = "#name#"
  description = "Service user that can access the dashboard"
}

resource "dynatrace_iam_group" "sample_group" {
  name        = "#name#"
  description = "Group that can acccess the dashboard"
}
