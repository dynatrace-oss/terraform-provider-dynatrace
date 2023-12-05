# ID vu9U3hXa3q0AAAABABlidWlsdGluOm9uZWFnZW50LmZlYXR1cmVzAAZ0ZW5hbnQABnRlbmFudAAkMWQzYjY4ODMtOWViZi0zMDljLTg1YjktNjg4OTcxYzE3NDM1vu9U3hXa3q0
resource "dynatrace_oneagent_side_masking" "#name#" {
  is_email_masking_enabled     = false
  is_financial_masking_enabled = true
  is_numbers_masking_enabled   = true
  is_query_masking_enabled     = false
}