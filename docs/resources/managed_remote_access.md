---
layout: ""
page_title: dynatrace_managed_remote_access Resource - terraform-provider-dynatrace"
subcategory: "Cluster Management"
description: |-
  The resource `dynatrace_managed_remote_access` covers configuration for remote access requests
---

# dynatrace_managed_remote_access (Resource)

!> **HTTP DELETE method not available** Terraform will no longer manage this resource on `destroy` but the configuration will still be present on the Dynatrace cluster.

-> This resource requires the cluster API token scope **Service Provider API** (`ServiceProviderAPI`)

## Dynatrace Documentation

- Cluster remote access - https://www.dynatrace.com/support/help/managed-cluster/configuration/cluster-remote-access

- Cluster API v2 - https://www.dynatrace.com/support/help/managed-cluster/cluster-api/cluster-api-v2

## Resource Example Usage

```terraform
resource "dynatrace_managed_remote_access" "Test" {
	user_id 		= "example@dynatrace.com"
	reason			= "Example"
	requested_days	= 1
	role			= "devops-admin"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `reason` (String) Request reason description, cannot be changed once created
- `requested_days` (Number) For how many days access is requested, cannot be changed once created
- `role` (String) Requested role, cannot be changed once created
- `user_id` (String) User id, cannot be changed once created

### Optional

- `state` (String) Access request state. Automatically set as `ACCEPTED` on create, state can be changed in subsequent updates.

### Read-Only

- `id` (String) The ID of this resource.
 