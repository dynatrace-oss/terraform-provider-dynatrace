resource "dynatrace_kubernetes_enrichment" "#name#" {
  scope = "environment"
  rules {
    rule {
      type    = "LABEL"
      source  = "#name#"
      target  = "dt.cost.product"
    }
    rule {
      type    = "LABEL"
      source  = "#name#"
      primary_grail_tag = true
    }
    rule {
      type    = "LABEL"
      source  = "#name#"
      target  = "dt.cost.product"
      primary_grail_tag = false
    }
    rule {
      type    = "ANNOTATION"
      source  = "#name#"
      target  = "dt.security_context"
    }
    rule {
      type    = "ANNOTATION"
      source  = "#name#"
      primary_grail_tag = true
    }
    rule {
      type    = "ANNOTATION"
      source  = "#name#"
      target  = "dt.security_context"
      primary_grail_tag = false
    }
  }
}