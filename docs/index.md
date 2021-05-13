# Dynatrace Provider

The Dynatrace provider is used to interact with the resources supported by the Dynatrace API. The provider needs to be configured with the proper credentials before it can be used.

Use the links to the left to learn about the available resources.

## Example
```
# Terraform 0.13+ uses the Terraform Registry:

terraform {
    required_version = "~> 0.13.0"
    required_providers {
        dynatrace = {
            source = "dynatrace.com/com/dynatrace"
        }
    }
}

# Configure the Dynatrace provider
provider "dynatrace" {
    dt_env_url    = "https://########.live.dynatrace.com"
    dt_api_token  = "################"
}

# Terraform 0.12- can be specified as:
terraform {
    required_providers {
        dynatrace = {
            version = "1.0.11"
        }
    }
}
```

## Schema
#### Optional
- **dt_api_token** (String, Sensitive) API Token with the correct permissions
- **dt_env_url** (String) The URL of your Dynatrace Environment

