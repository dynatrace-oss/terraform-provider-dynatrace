resource "dynatrace_document" "this" {
  type    = "dashboard"
  name    = "Example Dashboard"
  content = file(format("%s/example-dashboard.json", path.module))
}

data "dynatrace_documents" "all-dashboard-and-notebooks" {}
