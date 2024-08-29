resource "dynatrace_kubernetes_enrichment" "#name#" {
    scope = "environment"
  rules {
    rule {
      type    = "LABEL"
      enabled = true
      source  = "#name#"
      target  = "dt.cost.product"
    }
    rule {
      type    = "ANNOTATION"
      enabled = true
      source  = "#name#"
      target  = "dt.security_context"
    }
  }
}