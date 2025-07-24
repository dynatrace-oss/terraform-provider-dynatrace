data "dynatrace_iam_user" "user_a" {
  email = "a@example.com"
}

data "dynatrace_iam_user" "user_b" {
  email = "b@example.com"
}

data "dynatrace_iam_group" "example_group" {
  name = "Terraform Example"
}

resource "dynatrace_setting_permissions" "github_connection_access" {
  settings_object_id = dynatrace_github_connection.example_connection.id
  all_users = "none"
  users {
    user {
        user_id = data.dynatrace_iam_user.user_a.uid
        access = "write"
    }
    user {
        user_id = data.dynatrace_iam_user.user_b.uid
        access = "read"
    }
  }
  groups {
    group {
        group_id = data.dynatrace_iam_group.example_group.id
        access = "write"
    }
  }
}
