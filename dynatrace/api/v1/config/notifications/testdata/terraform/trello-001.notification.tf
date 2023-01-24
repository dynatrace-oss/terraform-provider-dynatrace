resource "dynatrace_notification" "#name#" {
  trello {
    name                = "#name#"
    description         = "trello-description"
    active              = false
    alerting_profile    = dynatrace_alerting_profile.Default.id
    application_key     = "trello-application-key"
    authorization_token = "#######"
    board_id            = "trello-board-id"
    list_id             = "trello-list-id"
    resolved_list_id    = "trello-resolved-list-id"
    text                = "trello-text"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
