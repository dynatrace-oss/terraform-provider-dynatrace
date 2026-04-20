resource "dynatrace_automation_approval" "wf_approval" {
  workflow_app_access_approval_enabled = true
  external_approvals_enabled = true
}
