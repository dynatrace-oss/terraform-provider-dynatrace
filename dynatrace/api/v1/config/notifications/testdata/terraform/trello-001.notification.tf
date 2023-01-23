resource "dynatrace_notification" "#name#" {
  trello {
    name                = "#name#"
    description         = "trello-description"
    active              = false
    alerting_profile    = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    application_key     = "trello-application-key"
    authorization_token = "#######"
    board_id            = "trello-board-id"
    list_id             = "trello-list-id"
    resolved_list_id    = "trello-resolved-list-id"
    text                = "trello-text"
  }
}