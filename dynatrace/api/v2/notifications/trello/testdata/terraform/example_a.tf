resource "dynatrace_trello_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active              = false
  name                = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile             = data.dynatrace_alerting_profile.Default.id
  application_key     = "trello-application-key"
  board_id            = "trello-board-id"
  list_id             = "trello-list-id"
  resolved_list_id    = "trello-resolved-list-id"
  text                = "trello-text"
  description         = "trello-description"
  authorization_token = "trello-authorization-token"
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}