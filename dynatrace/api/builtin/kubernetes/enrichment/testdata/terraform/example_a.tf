resource "dynatrace_kubernetes_enrichment" "#name#" {
    scope = "environment"
  rules {
    rule {
      type    = "LABEL"
      source  = "#name#"
      target  = "dt.cost.product"
    }
    rule {
      type    = "ANNOTATION"
      source  = "#name#"
      target  = "dt.security_context"
    }
  }
}