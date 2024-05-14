---
layout: ""
page_title: "dynatrace_iam_policies Data Source - terraform-provider-dynatrace"
subcategory: "IAM"
description: |-
  The data source `dynatrace_iam_policies` covers queries for policies available with the given credentials.
---

# dynatrace_iam_policies (Data Source)

You can use the attributes `environments`, `accounts` and `globals` to refine which policies you want to query for.
* The attribute `global` indicates whether the results should also contain global (Dynatrace defined) policies
* The attribute `environment` is an array of environment IDs.
* The results won't contain any environment specific policies if the attribute `environments` has been omitted
* The results will contain policies for all environments reachable via the given credentials if `environments` is set to `["*"]`
* The attribute `accounts` is an array of accounts UUIDs. Set this to `["*"]` if you want to receive account specific policies.
* The results won't contain any account specific policies if the attribute `accounts` has been omitted
## Example Usage

The following example queries for polices of all environments reachable via the given credentials, all accounts and all global policies.
```terraform
data "dynatrace_iam_policies" "all" {
  environments = ["*"]
  accounts     = ["*"]
  global       = true
}
```
The following example queries for policies that are defined for the environment with the id `abce234`. No account specific or global policies will be included.
```terraform
data "dynatrace_iam_policies" "all" {
  environments = ["abce234"]
  global       = false
}
```

## Example Output
```terraform
data "dynatrace_iam_policies" "all" {
  environments = ["*"]
  accounts     = ["*"]
  global       = true
}

output "policies" {
  value = data.dynatrace_iam_policies.all.policies
}
```

```
Changes to Outputs:
  + policies = [
      + {
          + account     = "########-86d8-####-88bd-############"
          + environment = ""
          + global      = false
          + id          = "########-7a6a-####-a43e-#############-#account#-#########-86d8-####-88bd-############"      
          + name        = "storage:bucket-definitions:delete"
          + uuid        = "########-7a6a-####-a43e-############"
        },
        ...
      + {
          + account     = ""
          + environment = "#######"
          + global      = false
          + id          = "########-c7d6-####-878c-#############-#environment#-########"
          + name        = "some-policy"
          + uuid        = "########-c7d6-####-878c-############"
        }, 
        ...
      + {
          + account     = ""
          + environment = ""
          + global      = true
          + id          = "########-6852-####-9d1b-#############-#global#-#global"
          + name        = "Storage Events Read"
          + uuid        = "########-6852-####-9d1b-############"
        },               
    ]

```
<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `accounts` (List of String) The results will contain policies defined for the given accountID. If one of the entries contains `*` the results will contain policies for all accounts
- `environments` (List of String) The results will contain policies defined for the given environments. If one of the entries contains `*` the results will contain policies for all environments
- `global` (Boolean) If `true` the results will contain global policies

### Read-Only

- `id` (String) The ID of this resource.
- `policies` (List of Object) (see [below for nested schema](#nestedatt--policies))

<a id="nestedatt--policies"></a>
### Nested Schema for `policies`

Read-Only:

- `account` (String)
- `environment` (String)
- `global` (Boolean)
- `id` (String)
- `name` (String)
- `uuid` (String)