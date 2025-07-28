data "dynatrace_iam_user" "user_a" {
  email = "a@example.com"
}

data "dynatrace_iam_user" "user_b" {
  email = "b@example.com"
}

data "dynatrace_iam_group" "example_group" {
  name = "Terraform Example"
}

resource "dynatrace_github_connection" "example_connection"{
  name    = "GitHub connection"
  type     = "pat"
  token   = "azAZ09"
}

resource "dynatrace_setting_permissions" "github_connection_access" {
  settings_object_id = dynatrace_github_connection.example_connection.id
  all_users = "none"
  users {
    user {
        uid = data.dynatrace_iam_user.user_a.uid
        access = "write"
    }
    user {
        uid = data.dynatrace_iam_user.user_b.uid
        access = "read"
    }
  }
  groups {
    group {
        id = data.dynatrace_iam_group.example_group.id
        access = "write"
    }
  }
}
