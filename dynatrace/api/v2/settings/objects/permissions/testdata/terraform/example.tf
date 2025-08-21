// make it more manual-execution friendly via variable
variable "IDENTIFIER" {
  type = string
  default = "#name#"
}

locals {
  // to use two users for one setting for read & write
  userCount = 2
  // to use two groups for one setting for read & write
  groupCount = 2
  // create 5 connections. One for users, one for groups, and one for all allUser variants
  // keep them separate so the combinations make sense (e.g., allUsers "write" with users/groups doesn't make sense)
  connectionCount = 5
}

resource "dynatrace_iam_group" "groups" {
  count = local.groupCount
  name = "${count.index}${var.IDENTIFIER}"
}

// for each group create a user
resource "dynatrace_iam_user" "users" {
  count = local.userCount
  email = "${var.IDENTIFIER}-${count.index}@example.com"
  groups = [dynatrace_iam_group.groups[count.index].id]
  # currently disabling, because for both example files we either have empty for both or non-empty for both and the other example should be visible in the documentation.
  # lifecycle {
  #   // currently the API adds a default group.
  #   // So there won't be the case that no updates are needed unless the default groups is added here
  #   // For the test to work, we have to ignore this one, else we'll get an error that terraform plan isn't empty
  #   ignore_changes = ["groups"]
  # }
}

// because the UID is not returned for the resource, we need data
data "dynatrace_iam_user" "users" {
  count = local.userCount
  email = dynatrace_iam_user.users[count.index].id
}

resource "dynatrace_github_connection" "example_connections" {
  count = local.connectionCount
  name    = "${count.index}-${var.IDENTIFIER}"
  type     = "pat"
  token   = "azAZ09"
}

resource "dynatrace_settings_permissions" "user_access" {
  settings_object_id = dynatrace_github_connection.example_connections[0].id
  all_users = "none"
  users {
    dynamic "user" {
      for_each = zipmap(data.dynatrace_iam_user.users[*].uid, ["read", "write"])
      content {
        uid = user.key
        access = user.value
      }
    }
  }
}

resource "dynatrace_settings_permissions" "group_access" {
  settings_object_id = dynatrace_github_connection.example_connections[1].id
  groups {
    dynamic "group" {
      for_each = zipmap(dynatrace_iam_group.groups[*].id, ["read", "write"])
      content {
        id = group.key
        access = group.value
      }
    }
  }
}

// Testing bug: for_each can't be used for a resource https://github.com/hashicorp/terraform-plugin-sdk/issues/536
resource "dynatrace_settings_permissions" "allUsers_access_none" {
  settings_object_id = dynatrace_github_connection.example_connections[2].id
  all_users          = "none"
}

resource "dynatrace_settings_permissions" "allUsers_access_read" {
  settings_object_id = dynatrace_github_connection.example_connections[3].id
  all_users          = "read"
}

resource "dynatrace_settings_permissions" "allUsers_access_write" {
  settings_object_id = dynatrace_github_connection.example_connections[4].id
  all_users          = "write"
}
