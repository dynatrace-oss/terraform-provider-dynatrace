resource "dynatrace_segment" "segment" {
  name      = "#name#"
  is_public = false
  includes {
    items {
      data_object = "_all_data_object"
      filter = jsonencode({
        "children" : [
          {
            "key" : {
              "range" : { "from" : 0, "to" : 16 },
              "textValue" : "k8s.cluster.name",
              "type" : "Key",
              "value" : "k8s.cluster.name"
            },
            "operator" : {
              "range" : { "from" : 17, "to" : 18 },
              "textValue" : "=",
              "type" : "ComparisonOperator",
              "value" : "="
            },
            "range" : { "from" : 0, "to" : 23 },
            "type" : "Statement",
            "value" : {
              "range" : { "from" : 19, "to" : 23 },
              "textValue" : "prod",
              "type" : "String",
              "value" : "prod"
            }
          }
        ],
        "explicit" : false,
        "logicalOperator" : "AND",
        "range" : { "from" : 0, "to" : 23 },
        "type" : "Group"
      })
    }
  }
}

data "dynatrace_segment" "segment" {
  name = dynatrace_segment.segment.name
}

data "dynatrace_segment" "by_id" {
  id = dynatrace_segment.segment.id
}

output "id" {
  value = data.dynatrace_segment.segment.id
}

output "name_by_id" {
  value = data.dynatrace_segment.by_id.name
}

output "is_public" {
  value = tostring(data.dynatrace_segment.segment.is_public)
}

data "dynatrace_segments" "all" {
  depends_on = [dynatrace_segment.segment]
}

output "segment_names" {
  value = join(",", data.dynatrace_segments.all.items[*].name)
}
