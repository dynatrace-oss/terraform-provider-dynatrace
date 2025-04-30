---
page_title: "dynatrace_dql Data Source - terraform-provider-dynatrace"
subcategory: "Live Data"
description: |-
  The data source `dynatrace_dql` allows for executing DQL queries. See detailed information about the [Dynatrace Query Language](https://docs.dynatrace.com/docs/discover-dynatrace/references/dynatrace-query-language) within the Dynatrace Documentation.
---

# dynatrace_dql (Data Source)

The only required attribute is the `query` attribute - holding the DQL query.
You may or may not utilize the additional attributes in order to narrow down the results. But most of that can also get achieved with in the DQL query itself.

The result of the query will be available within the attribute `results` in JSON format - usually an array of records.
The schema behind these results can, of course, vary depending on the DQL query you're executing.

Terraform will attempt to poll for results until the query has finished. There is no need to specify a timeout for that.

!> Executing DQL queries can inflict additional costs in Dynatrace. Be aware of that fact when using this Data Source. Terraform will run that query by default every time you're executing `terraform plan` or `terraform apply`.

## Example Usage
```terraform
data "dynatrace_dql" "this" {
  query = "fetch events"
}
```
will produce content for the `results` attribute like this:
```
[
    {
        "event.id"                             = "-7629786693326919096_1745910027748"
        "Event source"                         = "OS services monitoring"
        ...
        timestamp                              = "2025-04-29T07:00:44.416000000Z"
    },
    {
        ...
    },
    ...
]
```

You can also use Heredoc syntax for better readability of complex DQL queries.

```terraform
data "dynatrace_dql" "this" {
  query = <<EOT
    fetch events |
    filter event.type == "davis" AND davis.status != "CLOSED" |
    fields timestamp, davis.title, davis.underMaintenance, davis.status |
    sort timestamp |
    limit 10  
EOT
}
```

{{ .SchemaMarkdown | trimspace }}