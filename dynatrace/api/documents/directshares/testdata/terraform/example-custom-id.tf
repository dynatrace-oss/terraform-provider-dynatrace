resource "dynatrace_iam_group" "this" {
  name = "#name#"
}

resource "dynatrace_document" "this" {
  custom_id = "#name#"
  type      = "dashboard"
  name      = "#name#"
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

resource "dynatrace_direct_shares" "this" {
  document_id = dynatrace_document.this.id
  access      = "read"

  recipients {
    recipient {
      id   = dynatrace_iam_group.this.id
      type = "group"
    }
  }
}
