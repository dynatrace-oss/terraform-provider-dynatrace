resource "dynatrace_davis_copilot" "#name#" {
  enable_copilot                  = true
  enable_tenant_aware_data_mining = true
  blocklist_entries {
    blocklist_entrie {
      name = "#name#"
      type = "TABLE"
    }
  }
}
