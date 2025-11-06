---
layout: ""
page_title: "Resource - terraform-provider-dynatrace"
subcategory: "Grail"
description: |-
  The resource `dynatrace_segment` covers configuration of segments to logically structure and conveniently filter observability data across apps on the Dynatrace platform
---

# dynatrace_segment (Resource)

-> **Dynatrace SaaS only**

-> To utilize this resource with access to all segments, please define the environment variables `DT_CLIENT_ID`, `DT_CLIENT_SECRET`, `DT_ACCOUNT_ID` with an OAuth client including the following permissions: **View Filter-Segments** (`storage:filter-segments:read`), **Create and Update Filter-Segments** (`storage:filter-segments:write`), **Share Filter-Segments** (`storage:filter-segments:share`), **Delete Filter-Segments** (`storage:filter-segments:delete`) and **Maintain all Filter-Segments on the environment** (`storage:filter-segments:admin`).

-> This resource is excluded by default in the export utility, please explicitly specify the resource to retrieve existing configuration.

-> The `includes.items[X].filter` attribute, which is a JSON string, is unfriendly for configuration as code. The structure of that attribute is not publicly documented and therefore subject to change without warning. The resource schema has been created to match our REST API, but we will be reaching out to product management on further enhancement of this endpoint. In the meantime, please use the export utility to create configurations more efficiently.

## Dynatrace Documentation

- Segments - https://docs.dynatrace.com/docs/manage/segments

- Grail Storage Filter-Segments (API) - https://########.apps.dynatrace.com/platform/swagger-ui/index.html?urls.primaryName=Grail+-+Filter+Segments

## Resource Example Usage

```terraform
# ID GQ7NqJGPV1N
resource "dynatrace_segment" "#name#" {
  name        = "#name#"
  description = "Example description"
  is_public   = true
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
                                "textValue": "k8s.cluster.name",
                                "type": "Key",
                                "value": "k8s.cluster.name"
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
                                "to": 27
                          },
                          "type": "Statement",
                          "value": {
                                "range": {
                                      "from": 19,
                                      "to": 27
                                },
                                "textValue": "$cluster",
                                "type": "String",
                                "value": "$cluster"
                          }
                    }
              ],
              "explicit": false,
              "logicalOperator": "AND",
              "range": {
                    "from": 0,
                    "to": 27
              },
              "type": "Group"
        })
    }
    items {
      data_object = "dt.entity.cloud_application"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    items {
      data_object = "dt.entity.cloud_application_instance"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    items {
      data_object = "dt.entity.cloud_application_namespace"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    items {
      data_object = "dt.entity.container_group_instance"
      filter      = ""
      relationship {
        name   = "belongs_to"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    items {
      data_object = "dt.entity.host"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    items {
      data_object = "dt.entity.kubernetes_cluster"
      filter      = jsonencode({
              "children": [
                    {
                          "key": {
                                "range": {
                                      "from": 0,
                                      "to": 11
                                },
                                "textValue": "entity.name",
                                "type": "Key",
                                "value": "entity.name"
                          },
                          "operator": {
                                "range": {
                                      "from": 12,
                                      "to": 13
                                },
                                "textValue": "=",
                                "type": "ComparisonOperator",
                                "value": "="
                          },
                          "range": {
                                "from": 0,
                                "to": 22
                          },
                          "type": "Statement",
                          "value": {
                                "range": {
                                      "from": 14,
                                      "to": 22
                                },
                                "textValue": "$cluster",
                                "type": "String",
                                "value": "$cluster"
                          }
                    }
              ],
              "explicit": false,
              "logicalOperator": "AND",
              "range": {
                    "from": 0,
                    "to": 22
              },
              "type": "Group"
        })
    }
    items {
      data_object = "dt.entity.kubernetes_node"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    items {
      data_object = "dt.entity.kubernetes_service"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
    items {
      data_object = "dt.entity.service"
      filter      = ""
      relationship {
        name   = "clustered_by"
        target = "dt.entity.kubernetes_cluster"
      }
    }
  }
  variables {
    type  = "query"
    value =<<-EOT
      fetch dt.entity.kubernetes_cluster
      | fields cluster = entity.name
      | sort cluster
    EOT
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `is_public` (Boolean) Indicates if the filter-segment is publicly accessible within the tenant
- `name` (String) Name of the filter-segment

### Optional

- `description` (String) Description of the filter-segment
- `includes` (Block List, Max: 1) List of includes of the filter-segment (see [below for nested schema](#nestedblock--includes))
- `variables` (Block List, Max: 1) Variables of the filter-segment (see [below for nested schema](#nestedblock--variables))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--includes"></a>
### Nested Schema for `includes`

Optional:

- `items` (Block Set, Max: 20) TODO: No documentation available (see [below for nested schema](#nestedblock--includes--items))

<a id="nestedblock--includes--items"></a>
### Nested Schema for `includes.items`

Required:

- `data_object` (String) The data object that the filter will be applied to. Use '_all_data_object' to apply it to all dataObjects
- `filter` (String) Data will be filtered by this value

Optional:

- `apply_to` (Set of String) [Experimental] The tables that the entity-filter will be applied to`
- `relationship` (Block List, Max: 1) [Experimental] The relationship of an include which has to be be specified when the data object is an entity view (see [below for nested schema](#nestedblock--includes--items--relationship))

<a id="nestedblock--includes--items--relationship"></a>
### Nested Schema for `includes.items.relationship`

Required:

- `name` (String) Name of the relationship
- `target` (String) Target of the relationship




<a id="nestedblock--variables"></a>
### Nested Schema for `variables`

Required:

- `type` (String) Type of the variable
- `value` (String) Value of the variable
