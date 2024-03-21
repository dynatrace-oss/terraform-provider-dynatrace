resource "dynatrace_log_buckets" "first-instance" {
  enabled     = true
  bucket_name = "default_logs"
  matcher     = "matchesPhrase(content, \"error\")"
  rule_name   = "#name#"
}

resource "dynatrace_log_buckets" "second-instance" {
  enabled      = true
  bucket_name  = "default_logs"
  matcher      = "matchesPhrase(content, \"error\")"
  rule_name    = "#name#-second"
  insert_after = dynatrace_log_buckets.first-instance.id
}
