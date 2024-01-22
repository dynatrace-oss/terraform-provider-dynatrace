resource "dynatrace_log_buckets" "#name#" {
  enabled     = true
  bucket_name = "default_logs"
  matcher     = "matchesPhrase(content, \"error\")"
  rule_name   = "#name#"
}
