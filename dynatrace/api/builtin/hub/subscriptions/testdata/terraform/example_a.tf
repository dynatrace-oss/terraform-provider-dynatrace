resource "dynatrace_hub_subscriptions" "#name#" {
  token_subscriptions {
    token_subscription {
      name        = "#name"
      description = "Description"
      enabled     = true
      token       = "123456789012345678901234567890123456"
    }
  }
}
