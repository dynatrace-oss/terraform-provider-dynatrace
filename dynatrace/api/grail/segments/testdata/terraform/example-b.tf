# ID 1AL7hE2Ysvf
resource "dynatrace_segment" "#name#" {
  name        = "#name#"
  is_public   = false
  includes {
    items {
      data_object = "_all_data_object"
      filter      = jsonencode({
              "children": [
                    {
                          "key": {
                                "range": {
                                      "from": 0,
                                      "to": 16
                                },
                                "textValue": "dt.host_group.id",
                                "type": "Key",
                                "value": "dt.host_group.id"
                          },
                          "operator": {
                                "range": {
                                      "from": 17,
                                      "to": 18
                                },
                                "textValue": "=",
                                "type": "ComparisonOperator",
                                "value": "="
                          },
                          "range": {
                                "from": 0,
                                "to": 29
                          },
                          "type": "Statement",
                          "value": {
                                "range": {
                                      "from": 19,
                                      "to": 29
                                },
                                "textValue": "$hostgroup",
                                "type": "String",
                                "value": "$hostgroup"
                          }
                    }
              ],
              "explicit": false,
              "logicalOperator": "AND",
              "range": {
                    "from": 0,
                    "to": 29
              },
              "type": "Group"
        })
    }
    items {
      data_object = "dt.entity.host"
      filter      = jsonencode({
              "children": [
                    {
                          "key": {
                                "range": {
                                      "from": 0,
                                      "to": 13
                                },
                                "textValue": "hostGroupName",
                                "type": "Key",
                                "value": "hostGroupName"
                          },
                          "operator": {
                                "range": {
                                      "from": 14,
                                      "to": 15
                                },
                                "textValue": "=",
                                "type": "ComparisonOperator",
                                "value": "="
                          },
                          "range": {
                                "from": 0,
                                "to": 26
                          },
                          "type": "Statement",
                          "value": {
                                "range": {
                                      "from": 16,
                                      "to": 26
                                },
                                "textValue": "$hostgroup",
                                "type": "String",
                                "value": "$hostgroup"
                          }
                    }
              ],
              "explicit": false,
              "logicalOperator": "AND",
              "range": {
                    "from": 0,
                    "to": 26
              },
              "type": "Group"
        })
    }
    items {
      data_object = "dt.entity.host_group"
      filter      = ""
      relationship {
        name   = "instantiates"
        target = "dt.entity.host"
      }
    }
    items {
      data_object = "dt.entity.process_group"
      filter      = ""
      relationship {
        name   = "runs_on"
        target = "dt.entity.host"
      }
    }
    items {
      data_object = "dt.entity.process_group_instance"
      filter      = ""
      relationship {
        name   = "belongs_to"
        target = "dt.entity.host"
      }
    }
    items {
      data_object = "dt.entity.service"
      filter      = ""
      relationship {
        name   = "runs_on"
        target = "dt.entity.host"
      }
    }
  }
  variables {
    type  = "query"
    value =<<-EOT
      fetch dt.entity.host_group
      | fields hostgroup = entity.name
      | sort upper(hostgroup)
    EOT
  }
}
