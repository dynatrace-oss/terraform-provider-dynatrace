variable "random_name" {
  description = "A random name for use in the test"
  type        = string
}

resource "dynatrace_document" "sample_dashboard" {
  name = "${var.random_name}-dashboard"
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

resource "dynatrace_iam_group" "sample_group1" {
  name        = "${var.random_name}-group1"
  description = "First group that can acccess the dashboard"
}


resource "dynatrace_iam_group" "sample_group2" {
  name        = "${var.random_name}-group2"
  description = "Second group that can acccess the dashboard"
}

resource "dynatrace_direct_shares" "direct_share" {
  access      = "read-write"
  document_id = dynatrace_document.sample_dashboard.id
  recipients {
    recipient {
      type = "group"
      id   = dynatrace_iam_group.sample_group1.id
    }
    recipient {
      type = "group"
      id   = dynatrace_iam_group.sample_group2.id
    }
  }
}
