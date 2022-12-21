---
layout: ""
page_title: dynatrace_management_zone_v2 Resource - terraform-provider-dynatrace"
description: |-
  The resource `dynatrace_management_zone_v2` covers configuration for management zones
---

# dynatrace_management_zone_v2 (Resource)

## Dynatrace Documentation

- Management zones - https://www.dynatrace.com/support/help/how-to-use-dynatrace/management-zones

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:management-zones`)

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_management_zone_v2` downloads all existing management zone configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_management_zone_v2" "#name#" {
  name = "Cloud: Google"
  rules {
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION_NAMESPACE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "extensions"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type           = "HOST"
        host_to_pgpropagation = true
        attribute_conditions {
          condition {
            enum_value = "GOOGLE_CLOUD_PLATFORM"
            key        = "HOST_CLOUD_TYPE"
            operator   = "EQUALS"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CUSTOM_DEVICE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "CUSTOM_DEVICE_NAME"
            operator       = "CONTAINS"
            string_value   = "gcp"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION_NAMESPACE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "linkerd"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type                 = "SERVICE"
        service_to_host_propagation = true
        service_to_pgpropagation    = true
        attribute_conditions {
          condition {
            entity_id = "HOST_GROUP-8A29CED074001723"
            key       = "HOST_GROUP_ID"
            operator  = "EQUALS"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "KUBERNETES_CLUSTER"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "consul"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION_NAMESPACE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "gke"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "KUBERNETES_CLUSTER"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "extensions"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type           = "HOST"
        host_to_pgpropagation = true
        attribute_conditions {
          condition {
            entity_id = "HOST_GROUP-8A29CED074001723"
            key       = "HOST_GROUP_ID"
            operator  = "EQUALS"
          }
        }
      }
    }
    rule {
      type            = "DIMENSION"
      enabled         = true
      entity_selector = ""
      dimension_rule {
        applies_to = "METRIC"
        dimension_conditions {
          condition {
            condition_type = "METRIC_KEY"
            rule_matcher   = "BEGINS_WITH"
            value          = "cloud.gcp."
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "linkerd"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "WEB_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "WEB_APPLICATION_NAME"
            operator       = "CONTAINS"
            string_value   = "gcp"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type                 = "SERVICE"
        service_to_host_propagation = true
        service_to_pgpropagation    = true
        attribute_conditions {
          condition {
            enum_value = "GOOGLE_CLOUD_PLATFORM"
            key        = "HOST_CLOUD_TYPE"
            operator   = "EQUALS"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "extensions"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION_NAMESPACE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "consul"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "consul"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "gke"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "KUBERNETES_CLUSTER"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "linkerd"
          }
        }
      }
    }
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "WEB_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = true
            key            = "WEB_APPLICATION_NAME"
            operator       = "CONTAINS"
            string_value   = "www.gcp.easytravel.com"
          }
        }
      }
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) **Be careful when renaming** - if there are policies that are referencing this Management zone, they will need to be adapted to the new name!

### Optional

- `description` (String) Description
- `rules` (Block List, Max: 1) Rules (see [below for nested schema](#nestedblock--rules))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Optional:

- `rule` (Block Set) A management zone rule (see [below for nested schema](#nestedblock--rules--rule))

<a id="nestedblock--rules--rule"></a>
### Nested Schema for `rules.rule`

Required:

- `enabled` (Boolean) Enabled
- `type` (String) Rule type

Optional:

- `attribute_rule` (Block List, Max: 1) No documentation available (see [below for nested schema](#nestedblock--rules--rule--attribute_rule))
- `dimension_rule` (Block List, Max: 1) No documentation available (see [below for nested schema](#nestedblock--rules--rule--dimension_rule))
- `entity_selector` (String) Entity selector. The documentation of the entity selector can be found [here](https://dt-url.net/apientityselector).

<a id="nestedblock--rules--rule--attribute_rule"></a>
### Nested Schema for `rules.rule.attribute_rule`

Required:

- `attribute_conditions` (Block List, Min: 1, Max: 1) Conditions (see [below for nested schema](#nestedblock--rules--rule--attribute_rule--attribute_conditions))
- `entity_type` (String) Rule applies to

Optional:

- `azure_to_pgpropagation` (Boolean) Apply to process groups connected to matching Azure entities
- `azure_to_service_propagation` (Boolean) Apply to services provided by matching Azure entities
- `custom_device_group_to_custom_device_propagation` (Boolean) Apply to custom devices in a custom device group
- `host_to_pgpropagation` (Boolean) Apply to processes running on matching hosts
- `pg_to_host_propagation` (Boolean) Apply to underlying hosts of matching process groups
- `pg_to_service_propagation` (Boolean) Apply to all services provided by the process groups
- `service_to_host_propagation` (Boolean) Apply to underlying hosts of matching services
- `service_to_pgpropagation` (Boolean) Apply to underlying process groups of matching services

<a id="nestedblock--rules--rule--attribute_rule--attribute_conditions"></a>
### Nested Schema for `rules.rule.attribute_rule.attribute_conditions`

Optional:

- `condition` (Block Set) Attribute conditions (see [below for nested schema](#nestedblock--rules--rule--attribute_rule--attribute_conditions--condition))

<a id="nestedblock--rules--rule--attribute_rule--attribute_conditions--condition"></a>
### Nested Schema for `rules.rule.attribute_rule.attribute_conditions.condition`

Required:

- `key` (String) Property
- `operator` (String) Operator

Optional:

- `case_sensitive` (Boolean) Case sensitive
- `dynamic_key` (String) Dynamic key
- `dynamic_key_source` (String) Key source
- `entity_id` (String) Value
- `enum_value` (String) Value
- `integer_value` (Number) Value
- `string_value` (String) Value
- `tag` (String) Tag. Format: `[CONTEXT]tagKey:tagValue`




<a id="nestedblock--rules--rule--dimension_rule"></a>
### Nested Schema for `rules.rule.dimension_rule`

Required:

- `applies_to` (String) Type

Optional:

- `dimension_conditions` (Block List, Max: 1) Conditions (see [below for nested schema](#nestedblock--rules--rule--dimension_rule--dimension_conditions))

<a id="nestedblock--rules--rule--dimension_rule--dimension_conditions"></a>
### Nested Schema for `rules.rule.dimension_rule.dimension_conditions`

Optional:

- `condition` (Block Set) Dimension conditions (see [below for nested schema](#nestedblock--rules--rule--dimension_rule--dimension_conditions--condition))

<a id="nestedblock--rules--rule--dimension_rule--dimension_conditions--condition"></a>
### Nested Schema for `rules.rule.dimension_rule.dimension_conditions.condition`

Required:

- `condition_type` (String) Type
- `rule_matcher` (String) Operator
- `value` (String) Value

Optional:

- `key` (String) Key
 