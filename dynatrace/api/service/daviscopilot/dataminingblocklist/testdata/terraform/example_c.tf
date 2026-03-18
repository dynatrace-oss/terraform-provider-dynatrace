resource "dynatrace_davis_copilot" "agentic_ai" {
  enable_copilot                  = true
  enable_tenant_aware_data_mining = true
  enable_agentic_ai = true
  blocklist_entries {
    blocklist_entrie {
      name = "#name#"
      type = "TABLE"
    }
  }
}
