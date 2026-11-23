resource "dynatrace_ingest_enrichment_config" "example" {
  type         = "CUSTOM"
  value_source = "#name#"
  target       = "primary_tags.role"
  condition    = "isNotNull(os.architecture)"
}
