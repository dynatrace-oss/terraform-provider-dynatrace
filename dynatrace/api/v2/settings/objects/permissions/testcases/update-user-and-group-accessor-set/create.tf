resource "dynatrace_iam_group" "group" {
  count = 2
  name = "#name#-${count.index}"
}

resource "dynatrace_iam_user" "user" {
  count = 2
  email = "#name#${count.index}@example.com"
  groups = [dynatrace_iam_group.group[0].id]
  // default group added
  lifecycle {
    ignore_changes = [groups]
  }
}

// because the UID is not returned for the resource, we need data
data "dynatrace_iam_user" "user" {
  email = dynatrace_iam_user.user[0].id
}

data "dynatrace_iam_user" "user2" {
  email = dynatrace_iam_user.user[1].id
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
    # to update
    user {
      uid = data.dynatrace_iam_user.user2.uid
      access = "write"
    }
  }
  groups {
    group {
      id = dynatrace_iam_group.group[0].id
      access = "read"
    }
    # to update
    group {
      id = dynatrace_iam_group.group[1].id
      access = "read"
    }
  }
}
