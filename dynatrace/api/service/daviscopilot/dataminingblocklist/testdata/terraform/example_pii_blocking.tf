resource "dynatrace_davis_copilot" "pii_blocking" {
  enable_copilot       = true
  pii_blocking_enabled = true
  pii_blocking_types {
    canadian_social_insurance_number = true
    credit_card_number               = true
    email_address                    = true
    iban_bank_account                = true
    ip_address                       = true
    phone_number                     = true
    url_query_parameters             = true
    us_bank_number                   = true
    us_social_security_number        = true
  }
}
