# Getting Started with Terraform and the Dynatrace Provider

## Install Terraform CLI

Follow the simple installation steps in the official Terraform documentation [here](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli). 

Summary of steps:
1. Download binary
2. Set terraform binary available under PATH
3. Verify installation by opening a new terminal and executing `terraform -help`

## Retrieve the Dynatrace Terraform Provider

The Dynatrace Provider is available in the Terraform Registry and can be automatically downloaded when initializing a working directory with `terraform init`. 

To get started, create a working directory with a `providers.tf` file with the example configuration block below. Replace version `x.y.z` with the latest release available of the provider available on [GitHub](https://github.com/dynatrace-oss/terraform-provider-dynatrace), e.g. 1.34.0

Example:
```terraform
terraform {
    required_providers {
        dynatrace = {
            version = "x.y.z"
            source = "dynatrace-oss/dynatrace"
        }
    }
} 
```

Open a terminal window and navigate to the working directory; run `terraform init`.

```terraform
Initializing the backend...

Initializing provider plugins...
- Finding dynatrace-oss/dynatrace versions matching "x.y.z"...
- Installing dynatrace-oss/dynatrace x.y.z...
- Installed dynatrace-oss/dynatrace x.y.z (signed by a HashiCorp partner, key ID *************)
...

Terraform has been successfully initialized!
```

## What's Next
Hands-on: [Terraform Dynatrace Basic Example](LINK)