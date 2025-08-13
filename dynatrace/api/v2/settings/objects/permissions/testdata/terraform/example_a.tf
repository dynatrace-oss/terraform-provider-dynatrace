resource "dynatrace_iam_group" "group" {
  name = "#name#"
}

// for each group create a user
resource "dynatrace_iam_user" "user" {
  email = "#name#@example.com"
  groups = [dynatrace_iam_group.group.id]
}

// because the UID is not returned for the resource, we need data
data "dynatrace_iam_user" "user" {
  email = dynatrace_iam_user.user.id
}

resource "dynatrace_github_connection" "connection" {
  name    = "#name#"
  type     = "pat"
  token   = "azAZ09"
}

resource "dynatrace_settings_permissions" "permission" {
  settings_object_id = dynatrace_github_connection.connection.id
  all_users = "none"
  users {
    user {
      uid = data.dynatrace_iam_user.user.uid
      access = "write"
    }
  }
  groups {
    group {
      id = dynatrace_iam_group.group.id
      access = "read"
    }
  }
}