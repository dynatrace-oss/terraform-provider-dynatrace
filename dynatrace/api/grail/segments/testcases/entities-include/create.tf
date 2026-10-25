resource "dynatrace_segment" "segment" {
  name      = "#name#"
  is_public = false
  variables {
    type  = "query"
    value = <<-EOT
      metrics
       | filter isNotNull(primary_tags.team) and primary_tags.team == "payments"
       | fields primary_tags.team, dt.smartscape.service
       | summarize by: { primary_tags.team }, dt.smartscape.service = collectDistinct(dt.smartscape.service)
       | sort lower(primary_tags.team)
    EOT
  }
  includes {
    items {
      data_object = "_all_data_object"
      filter = jsonencode({
        type                    = "Group"
        range                   = { from = 0, to = 39 }
        logicalOperator         = "AND"
        explicit                = false
        explicitLogicalOperator = false
        version                 = "007"
        children = [{
          type  = "Statement"
          range = { from = 0, to = 38 }
          key = {
            type      = "Key"
            range     = { from = 0, to = 17 }
            textValue = "primary_tags.team"
            value     = "primary_tags.team"
          }
          operator = {
            type      = "ComparisonOperator"
            range     = { from = 18, to = 19 }
            textValue = "="
            value     = "="
          }
          value = {
            type         = "Variable"
            range        = { from = 20, to = 38 }
            textValue    = "$primary_tags.team"
            value        = "$primary_tags.team"
            variableType = "Any"
          }
        }]
      })
    }
    items {
      data_object = "_all_entities"
      filter = jsonencode({
        type                    = "Group"
        range                   = { from = 0, to = 51 }
        logicalOperator         = "AND"
        explicit                = false
        explicitLogicalOperator = true
        version                 = "007"
        children = [
          {
            type  = "Statement"
            range = { from = 0, to = 16 }
            key = {
              type      = "Key"
              range     = { from = 0, to = 4 }
              textValue = "type"
              value     = "type"
            }
            operator = {
              type      = "ComparisonOperator"
              range     = { from = 5, to = 6 }
              textValue = "="
              value     = "="
            }
            value = {
              type      = "String"
              range     = { from = 7, to = 16 }
              textValue = "\"SERVICE\""
              value     = "SERVICE"
            }
          },
          {
            type  = "Statement"
            range = { from = 21, to = 51 }
            key = {
              type      = "Key"
              range     = { from = 21, to = 23 }
              textValue = "id"
              value     = "id"
            }
            operator = {
              type      = "InclusionOperator"
              range     = { from = 24, to = 26 }
              textValue = "in"
              value     = "in"
            }
            value = {
              type      = "List"
              range     = { from = 27, to = 51 }
              textValue = "($dt.smartscape.service)"
              value = [{
                type         = "Variable"
                range        = { from = 28, to = 50 }
                textValue    = "$dt.smartscape.service"
                value        = "$dt.smartscape.service"
                variableType = "Any"
              }]
            }
          }
        ]
      })
    }
  }
}
