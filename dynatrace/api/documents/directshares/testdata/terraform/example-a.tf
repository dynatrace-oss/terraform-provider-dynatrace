resource "dynatrace_direct_shares" "this" {
  document_id = "441564f0-23c9-40ef-b344-18c02c23d712"
  access      = "read-write"

  recipients {
    recipient {
      id   = "441664f0-23c9-40ef-b344-18c02c23d787"
      type = "user"
    }

    recipient {
      id   = "441664f0-23c9-40ef-b344-18c02c23d788"
      type = "group"
    }
  }
}
