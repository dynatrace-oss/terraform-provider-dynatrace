---
layout: ""
page_title: dynatrace_management_zone_v2 Resource - terraform-provider-dynatrace"
subcategory: "Management Zones"
description: |-
  The resource `dynatrace_management_zone_v2` covers configuration for management zones
---

# dynatrace_management_zone_v2 (Resource)

-> This resource requires the API token scopes **Read settings** (`settings.read`) and **Write settings** (`settings.write`)

## Dynatrace Documentation

- Management zones - https://www.dynatrace.com/support/help/how-to-use-dynatrace/management-zones

- Settings API - https://www.dynatrace.com/support/help/dynatrace-api/environment-api/settings (schemaId: `builtin:management-zones`)

## Environment Variables (Optional)

There may be a delay for this resource to be fully available as a dependency for a subsequent resource. E.g. Utilizing this resource and `dynatrace_slo` together.
 
A default polling mechanism exists to validate the creation but may require adjustment due to load. The following environment variables can be used to fine tune these settings.

- `DT_MGMZ_RETRIES` (Default: 50, Max: 150) configures the maximum attempts to confirm that the create operation has succeeded.
- `DT_MGMZ_SUCCESSES` (Default: 5, Max: 25) configures the number of successful consecutive retries expected.

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_management_zone_v2` downloads all existing management zone configuration

The full documentation of the export feature is available [here](https://dt-url.net/h203qmc).

## Resource Example Usage

```terraform
resource "dynatrace_management_zone_v2" "#name#" {
  name = "#name#"
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
- `legacy_id` (String) The ID of this setting when referred to by the Config REST API V1
- `rules` (Block List, Max: 1) Rules (see [below for nested schema](#nestedblock--rules))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Required:

- `rule` (Block Set, Min: 1) A management zone rule (see [below for nested schema](#nestedblock--rules--rule))

<a id="nestedblock--rules--rule"></a>
### Nested Schema for `rules.rule`

Required:

- `enabled` (Boolean) This setting is enabled (`true`) or disabled (`false`)
- `type` (String) Possible Values: `DIMENSION`, `ME`, `SELECTOR`

Optional:

- `attribute_rule` (Block List, Max: 1) no documentation available (see [below for nested schema](#nestedblock--rules--rule--attribute_rule))
- `dimension_rule` (Block List, Max: 1) no documentation available (see [below for nested schema](#nestedblock--rules--rule--dimension_rule))
- `entity_selector` (String) The documentation of the entity selector can be found [here](https://dt-url.net/apientityselector).

<a id="nestedblock--rules--rule--attribute_rule"></a>
### Nested Schema for `rules.rule.attribute_rule`

Required:

- `attribute_conditions` (Block List, Min: 1, Max: 1) no documentation available (see [below for nested schema](#nestedblock--rules--rule--attribute_rule--attribute_conditions))
- `entity_type` (String) Possible Values: `APPMON_SERVER`, `APPMON_SYSTEM_PROFILE`, `AWS_ACCOUNT`, `AWS_APPLICATION_LOAD_BALANCER`, `AWS_AUTO_SCALING_GROUP`, `AWS_CLASSIC_LOAD_BALANCER`, `AWS_NETWORK_LOAD_BALANCER`, `AWS_RELATIONAL_DATABASE_SERVICE`, `AZURE`, `BROWSER_MONITOR`, `CLOUD_APPLICATION`, `CLOUD_APPLICATION_NAMESPACE`, `CLOUD_FOUNDRY_FOUNDATION`, `CUSTOM_APPLICATION`, `CUSTOM_DEVICE`, `CUSTOM_DEVICE_GROUP`, `DATA_CENTER_SERVICE`, `ENTERPRISE_APPLICATION`, `ESXI_HOST`, `EXTERNAL_MONITOR`, `HOST`, `HOST_GROUP`, `HTTP_MONITOR`, `KUBERNETES_CLUSTER`, `KUBERNETES_SERVICE`, `MOBILE_APPLICATION`, `OPENSTACK_ACCOUNT`, `PROCESS_GROUP`, `QUEUE`, `SERVICE`, `WEB_APPLICATION`

Optional:

- `azure_to_pgpropagation` (Boolean) Apply to process groups connected to matching Azure entities
- `azure_to_service_propagation` (Boolean) Apply to services provided by matching Azure entities
- `custom_device_group_to_custom_device_propagation` (Boolean) Apply to custom devices in a custom device group
- `host_to_pgpropagation` (Boolean) Apply to processes running on matching hosts. `entity_type` must be set to `HOST`
- `pg_to_host_propagation` (Boolean) Apply to underlying hosts of matching process groups. `entity_type` must be set to `PROCESS_GROUP`
- `pg_to_service_propagation` (Boolean) Apply to all services provided by the process groups. `entity_type` must be set to `PROCESS_GROUP`
- `service_to_host_propagation` (Boolean) Apply to underlying hosts of matching services. `entity_type` must be set to `SERVICE`
- `service_to_pgpropagation` (Boolean) Apply to underlying process groups of matching services. `entity_type` must be set to `SERVICE`

<a id="nestedblock--rules--rule--attribute_rule--attribute_conditions"></a>
### Nested Schema for `rules.rule.attribute_rule.attribute_conditions`

Required:

- `condition` (Block Set, Min: 1) Attribute conditions (see [below for nested schema](#nestedblock--rules--rule--attribute_rule--attribute_conditions--condition))

<a id="nestedblock--rules--rule--attribute_rule--attribute_conditions--condition"></a>
### Nested Schema for `rules.rule.attribute_rule.attribute_conditions.condition`

Required:

- `key` (String) Possible Values: `APPMON_SERVER_NAME`, `APPMON_SYSTEM_PROFILE_NAME`, `AWS_ACCOUNT_ID`, `AWS_ACCOUNT_NAME`, `AWS_APPLICATION_LOAD_BALANCER_NAME`, `AWS_APPLICATION_LOAD_BALANCER_TAGS`, `AWS_AUTO_SCALING_GROUP_NAME`, `AWS_AUTO_SCALING_GROUP_TAGS`, `AWS_AVAILABILITY_ZONE_NAME`, `AWS_CLASSIC_LOAD_BALANCER_FRONTEND_PORTS`, `AWS_CLASSIC_LOAD_BALANCER_NAME`, `AWS_CLASSIC_LOAD_BALANCER_TAGS`, `AWS_NETWORK_LOAD_BALANCER_NAME`, `AWS_NETWORK_LOAD_BALANCER_TAGS`, `AWS_RELATIONAL_DATABASE_SERVICE_DB_NAME`, `AWS_RELATIONAL_DATABASE_SERVICE_ENDPOINT`, `AWS_RELATIONAL_DATABASE_SERVICE_ENGINE`, `AWS_RELATIONAL_DATABASE_SERVICE_INSTANCE_CLASS`, `AWS_RELATIONAL_DATABASE_SERVICE_NAME`, `AWS_RELATIONAL_DATABASE_SERVICE_PORT`, `AWS_RELATIONAL_DATABASE_SERVICE_TAGS`, `AZURE_ENTITY_NAME`, `AZURE_ENTITY_TAGS`, `AZURE_MGMT_GROUP_NAME`, `AZURE_MGMT_GROUP_UUID`, `AZURE_REGION_NAME`, `AZURE_SCALE_SET_NAME`, `AZURE_SUBSCRIPTION_NAME`, `AZURE_SUBSCRIPTION_UUID`, `AZURE_TENANT_NAME`, `AZURE_TENANT_UUID`, `AZURE_VM_NAME`, `BROWSER_MONITOR_NAME`, `BROWSER_MONITOR_TAGS`, `CLOUD_APPLICATION_LABELS`, `CLOUD_APPLICATION_NAME`, `CLOUD_APPLICATION_NAMESPACE_LABELS`, `CLOUD_APPLICATION_NAMESPACE_NAME`, `CLOUD_FOUNDRY_FOUNDATION_NAME`, `CLOUD_FOUNDRY_ORG_NAME`, `CUSTOM_APPLICATION_NAME`, `CUSTOM_APPLICATION_PLATFORM`, `CUSTOM_APPLICATION_TAGS`, `CUSTOM_APPLICATION_TYPE`, `CUSTOM_DEVICE_DNS_ADDRESS`, `CUSTOM_DEVICE_GROUP_NAME`, `CUSTOM_DEVICE_GROUP_TAGS`, `CUSTOM_DEVICE_IP_ADDRESS`, `CUSTOM_DEVICE_METADATA`, `CUSTOM_DEVICE_NAME`, `CUSTOM_DEVICE_PORT`, `CUSTOM_DEVICE_TAGS`, `CUSTOM_DEVICE_TECHNOLOGY`, `DATA_CENTER_SERVICE_DECODER_TYPE`, `DATA_CENTER_SERVICE_IP_ADDRESS`, `DATA_CENTER_SERVICE_METADATA`, `DATA_CENTER_SERVICE_NAME`, `DATA_CENTER_SERVICE_PORT`, `DATA_CENTER_SERVICE_TAGS`, `DOCKER_CONTAINER_NAME`, `DOCKER_FULL_IMAGE_NAME`, `DOCKER_IMAGE_VERSION`, `EC2_INSTANCE_AMI_ID`, `EC2_INSTANCE_AWS_INSTANCE_TYPE`, `EC2_INSTANCE_AWS_SECURITY_GROUP`, `EC2_INSTANCE_BEANSTALK_ENV_NAME`, `EC2_INSTANCE_ID`, `EC2_INSTANCE_NAME`, `EC2_INSTANCE_PRIVATE_HOST_NAME`, `EC2_INSTANCE_PUBLIC_HOST_NAME`, `EC2_INSTANCE_TAGS`, `ENTERPRISE_APPLICATION_DECODER_TYPE`, `ENTERPRISE_APPLICATION_IP_ADDRESS`, `ENTERPRISE_APPLICATION_METADATA`, `ENTERPRISE_APPLICATION_NAME`, `ENTERPRISE_APPLICATION_PORT`, `ENTERPRISE_APPLICATION_TAGS`, `ESXI_HOST_CLUSTER_NAME`, `ESXI_HOST_HARDWARE_MODEL`, `ESXI_HOST_HARDWARE_VENDOR`, `ESXI_HOST_NAME`, `ESXI_HOST_PRODUCT_NAME`, `ESXI_HOST_PRODUCT_VERSION`, `ESXI_HOST_TAGS`, `EXTERNAL_MONITOR_ENGINE_DESCRIPTION`, `EXTERNAL_MONITOR_ENGINE_NAME`, `EXTERNAL_MONITOR_ENGINE_TYPE`, `EXTERNAL_MONITOR_NAME`, `EXTERNAL_MONITOR_TAGS`, `GEOLOCATION_SITE_NAME`, `GOOGLE_CLOUD_PLATFORM_ZONE_NAME`, `GOOGLE_COMPUTE_INSTANCE_ID`, `GOOGLE_COMPUTE_INSTANCE_MACHINE_TYPE`, `GOOGLE_COMPUTE_INSTANCE_NAME`, `GOOGLE_COMPUTE_INSTANCE_PROJECT`, `GOOGLE_COMPUTE_INSTANCE_PROJECT_ID`, `GOOGLE_COMPUTE_INSTANCE_PUBLIC_IP_ADDRESSES`, `HOST_AIX_LOGICAL_CPU_COUNT`, `HOST_AIX_SIMULTANEOUS_THREADS`, `HOST_AIX_VIRTUAL_CPU_COUNT`, `HOST_ARCHITECTURE`, `HOST_AWS_NAME_TAG`, `HOST_AZURE_COMPUTE_MODE`, `HOST_AZURE_SKU`, `HOST_AZURE_WEB_APPLICATION_HOST_NAMES`, `HOST_AZURE_WEB_APPLICATION_SITE_NAMES`, `HOST_BITNESS`, `HOST_BOSH_AVAILABILITY_ZONE`, `HOST_BOSH_DEPLOYMENT_ID`, `HOST_BOSH_INSTANCE_ID`, `HOST_BOSH_INSTANCE_NAME`, `HOST_BOSH_NAME`, `HOST_BOSH_STEMCELL_VERSION`, `HOST_CLOUD_TYPE`, `HOST_CPU_CORES`, `HOST_CUSTOM_METADATA`, `HOST_DETECTED_NAME`, `HOST_GROUP_ID`, `HOST_GROUP_NAME`, `HOST_HYPERVISOR_TYPE`, `HOST_IP_ADDRESS`, `HOST_KUBERNETES_LABELS`, `HOST_LOGICAL_CPU_CORES`, `HOST_NAME`, `HOST_ONEAGENT_CUSTOM_HOST_NAME`, `HOST_OS_TYPE`, `HOST_OS_VERSION`, `HOST_PAAS_MEMORY_LIMIT`, `HOST_PAAS_TYPE`, `HOST_TAGS`, `HOST_TECHNOLOGY`, `HTTP_MONITOR_NAME`, `HTTP_MONITOR_TAGS`, `KUBERNETES_CLUSTER_NAME`, `KUBERNETES_NODE_NAME`, `KUBERNETES_SERVICE_NAME`, `MOBILE_APPLICATION_NAME`, `MOBILE_APPLICATION_PLATFORM`, `MOBILE_APPLICATION_TAGS`, `NAME_OF_COMPUTE_NODE`, `OPENSTACK_ACCOUNT_NAME`, `OPENSTACK_ACCOUNT_PROJECT_NAME`, `OPENSTACK_AVAILABILITY_ZONE_NAME`, `OPENSTACK_PROJECT_NAME`, `OPENSTACK_REGION_NAME`, `OPENSTACK_VM_INSTANCE_TYPE`, `OPENSTACK_VM_NAME`, `OPENSTACK_VM_SECURITY_GROUP`, `PROCESS_GROUP_AZURE_HOST_NAME`, `PROCESS_GROUP_AZURE_SITE_NAME`, `PROCESS_GROUP_CUSTOM_METADATA`, `PROCESS_GROUP_DETECTED_NAME`, `PROCESS_GROUP_ID`, `PROCESS_GROUP_LISTEN_PORT`, `PROCESS_GROUP_NAME`, `PROCESS_GROUP_PREDEFINED_METADATA`, `PROCESS_GROUP_TAGS`, `PROCESS_GROUP_TECHNOLOGY`, `PROCESS_GROUP_TECHNOLOGY_EDITION`, `PROCESS_GROUP_TECHNOLOGY_VERSION`, `QUEUE_NAME`, `QUEUE_TECHNOLOGY`, `QUEUE_VENDOR`, `SERVICE_AKKA_ACTOR_SYSTEM`, `SERVICE_CTG_SERVICE_NAME`, `SERVICE_DATABASE_HOST_NAME`, `SERVICE_DATABASE_NAME`, `SERVICE_DATABASE_TOPOLOGY`, `SERVICE_DATABASE_VENDOR`, `SERVICE_DETECTED_NAME`, `SERVICE_ESB_APPLICATION_NAME`, `SERVICE_IBM_CTG_GATEWAY_URL`, `SERVICE_MESSAGING_LISTENER_CLASS_NAME`, `SERVICE_NAME`, `SERVICE_PORT`, `SERVICE_PUBLIC_DOMAIN_NAME`, `SERVICE_REMOTE_ENDPOINT`, `SERVICE_REMOTE_SERVICE_NAME`, `SERVICE_TAGS`, `SERVICE_TECHNOLOGY`, `SERVICE_TECHNOLOGY_EDITION`, `SERVICE_TECHNOLOGY_VERSION`, `SERVICE_TOPOLOGY`, `SERVICE_TYPE`, `SERVICE_WEB_APPLICATION_ID`, `SERVICE_WEB_CONTEXT_ROOT`, `SERVICE_WEB_SERVER_ENDPOINT`, `SERVICE_WEB_SERVER_NAME`, `SERVICE_WEB_SERVICE_NAME`, `SERVICE_WEB_SERVICE_NAMESPACE`, `VMWARE_DATACENTER_NAME`, `VMWARE_VM_NAME`, `WEB_APPLICATION_NAME`, `WEB_APPLICATION_NAME_PATTERN`, `WEB_APPLICATION_TAGS`, `WEB_APPLICATION_TYPE`
- `operator` (String) Possible Values: `BEGINS_WITH`, `CONTAINS`, `ENDS_WITH`, `EQUALS`, `EXISTS`, `GREATER_THAN`, `GREATER_THAN_OR_EQUAL`, `IS_IP_IN_RANGE`, `LOWER_THAN`, `LOWER_THAN_OR_EQUAL`, `NOT_BEGINS_WITH`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUALS`, `NOT_EXISTS`, `NOT_GREATER_THAN`, `NOT_GREATER_THAN_OR_EQUAL`, `NOT_IS_IP_IN_RANGE`, `NOT_LOWER_THAN`, `NOT_LOWER_THAN_OR_EQUAL`, `NOT_REGEX_MATCHES`, `NOT_TAG_KEY_EQUALS`, `REGEX_MATCHES`, `TAG_KEY_EQUALS`

Optional:

- `case_sensitive` (Boolean) Case sensitive
- `dynamic_key` (String) Dynamic key
- `dynamic_key_source` (String) Key source
- `entity_id` (String) Value
- `enum_value` (String) Value
- `integer_value` (Number) Value
- `string_value` (String) Value
- `tag` (String) Format: `[CONTEXT]tagKey:tagValue`




<a id="nestedblock--rules--rule--dimension_rule"></a>
### Nested Schema for `rules.rule.dimension_rule`

Required:

- `applies_to` (String) Possible Values: `ANY`, `LOG`, `METRIC`

Optional:

- `dimension_conditions` (Block List, Max: 1) Conditions (see [below for nested schema](#nestedblock--rules--rule--dimension_rule--dimension_conditions))

<a id="nestedblock--rules--rule--dimension_rule--dimension_conditions"></a>
### Nested Schema for `rules.rule.dimension_rule.dimension_conditions`

Required:

- `condition` (Block Set, Min: 1) Dimension conditions (see [below for nested schema](#nestedblock--rules--rule--dimension_rule--dimension_conditions--condition))

<a id="nestedblock--rules--rule--dimension_rule--dimension_conditions--condition"></a>
### Nested Schema for `rules.rule.dimension_rule.dimension_conditions.condition`

Required:

- `condition_type` (String) Possible Values: `DIMENSION`, `LOG_FILE_NAME`, `METRIC_KEY`
- `rule_matcher` (String) Possible Values: `BEGINS_WITH`, `EQUALS`
- `value` (String) no documentation available

Optional:

- `key` (String) no documentation available
 