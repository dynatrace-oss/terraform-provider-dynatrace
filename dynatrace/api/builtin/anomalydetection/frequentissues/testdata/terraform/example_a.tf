resource "dynatrace_frequent_issues" "#name#" {
  detect_apps = true
  detect_txn = true
  detect_infra = true
}
