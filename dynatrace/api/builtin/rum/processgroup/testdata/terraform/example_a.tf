variable "PROCESS_GROUP_ID" {
  type = string
}


resource "dynatrace_process_group_rum" "rum" {
  enable           = false
  process_group_id = var.PROCESS_GROUP_ID
}
