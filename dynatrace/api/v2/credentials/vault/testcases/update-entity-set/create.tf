resource "dynatrace_iam_service_user" "user" {
  count = 3
  name = "#name#"
}

resource "dynatrace_credentials" "cred" {
  name              = "#name#"
  username          = "username"
  password          = "password"
  scopes            = ["SYNTHETIC"]
  allowed_entities {
    entity {
      id = dynatrace_iam_service_user.user[0].id
      type = "USER"
    }
    # to update
    entity {
      id = dynatrace_iam_service_user.user[1].id
      type = "USER"
    }
  }
}
