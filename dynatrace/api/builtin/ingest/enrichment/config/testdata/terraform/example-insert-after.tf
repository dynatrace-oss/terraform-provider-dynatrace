resource "dynatrace_ingest_enrichment_config" "first" {
  type         = "CUSTOM"
  value_source = "#name#-1"
  target       = "primary_tags.role"
  condition    = "isNotNull(os.architecture)"
}

resource "dynatrace_ingest_enrichment_config" "second" {
  type         = "CUSTOM"
  value_source = "#name#-2"
  target       = "primary_tags.role"
  condition    = "isNotNull(os.architecture)"
  insert_after = dynatrace_ingest_enrichment_config.first.id
}
