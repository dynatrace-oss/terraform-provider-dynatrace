# Implementation Plan: IAM Service User Resource and Datasource

## Executive Summary
Implement `dynatrace_iam_service_user` resource and datasource to manage service users in Dynatrace IAM. Service users are non-human identities used for API access and automation. This implementation will follow the same patterns as IAM users and IAM groups.

## API Overview
Based on the Dynatrace API spec (https://api.dynatrace.com/spec/#/Service%20user%20management):

**Endpoints**:
- `POST /iam/v1/accounts/{accountUuid}/service-users` - Create service user
- `GET /iam/v1/accounts/{accountUuid}/service-users/{serviceUserId}` - Get service user
- `PUT /iam/v1/accounts/{accountUuid}/service-users/{serviceUserId}` - Update service user
- `DELETE /iam/v1/accounts/{accountUuid}/service-users/{serviceUserId}` - Delete service user
- `GET /iam/v1/accounts/{accountUuid}/service-users` - List service users. Results are returned in pages.

**Service User Properties**:
- `uid` (UUID) - Service user identifier (read-only)
- `email` (string) - Service user email (read-only)
- `name` (string) - Service user name (required)
- `description` (string) - Service user description (optional)
- `groups` (array of UUIDs) - Group memberships (optional)

## Current State
**Does Not Exist**: No implementation currently exists for IAM service users in the codebase.

## Implementation Phases

### Phase 1: API Service Implementation
**Location**: `dynatrace/api/iam/serviceusers/` (new directory)

#### 1.1 Service Client
**File**: `dynatrace/api/iam/serviceusers/service.go` (new)

Implement service client following the pattern from `dynatrace/api/iam/groups/service.go`:

```go
type ServiceUserServiceClient struct {
    clientID     string
    accountID    string
    clientSecret string
    tokenURL     string
    endpointURL  string
}

// CRUD operations:
func (me *ServiceUserServiceClient) Create(ctx context.Context, serviceUser *settings.ServiceUser) (*api.Stub, error)
func (me *ServiceUserServiceClient) Get(ctx context.Context, id string, v *settings.ServiceUser) error
func (me *ServiceUserServiceClient) Update(ctx context.Context, id string, serviceUser *settings.ServiceUser) error
func (me *ServiceUserServiceClient) Delete(ctx context.Context, id string) error
func (me *ServiceUserServiceClient) List(ctx context.Context) (api.Stubs, error)
func (me *ServiceUserServiceClient) SchemaID() string // "accounts:iam:serviceusers"
```

**Key Implementation Details**:

- Base URL: `{endpointURL}/iam/v1/accounts/{accountUuid}/service-users`
- Use IAM OAuth client `iam.NewIAMClient(me)`
- Strip `urn:dtaccount:` prefix from account ID
- Handle separate group assignment `PUT [/service-users/{id}/groups`
- Return UUID as resource ID
- Implement caching with revision tracking (like groups service)

#### 1.2 Settings Model

**File**: `dynatrace/api/iam/serviceusers/settings/service_user.go` (new)

```go
type ServiceUser struct {
    ID          string   `json:"id,omitempty"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Groups      []string `json:"-"`
}

func (me *ServiceUser) Schema() map[string]*schema.Schema
func (me *ServiceUser) MarshalHCL(properties hcl.Properties) error
func (me *ServiceUser) UnmarshalHCL(decoder hcl.Decoder) error
```

**Schema Fields**:

- `name` (Required, string) - Service user name
- `description` (Optional, string) - Description
- `groups` (Optional, Set of strings) - Group UUIDs
- `id` (Computed, string) - Service user UUID

#### 1.3 Provider Registration

**File**: `provider/provider.go`

Add imports and registration:

```go
import (
    rs_iamserviceuser "github.com/dynatrace-oss/terraform-provider-dynatrace/resources/iamserviceuser"
)

// In ResourcesMap:
"dynatrace_iam_service_user": rs_iamserviceuser.Resource(),
```

### Phase 2: Datasource Implementation

**Location**: `datasources/iam/serviceusers/` (new directory)

#### 2.1 Single Service User Datasource

**File**: `data_source.go` (new)

Follow pattern from `datasources/iam/users/data_source.go`

```go
func DataSource() *schema.Resource {
    return &schema.Resource{
        ReadContext: logging.EnableDSCtx(DataSourceRead),
        Description: "Fetches IAM service user details by name or ID",
        Schema: map[string]*schema.Schema{
            "name": {
                Type:         schema.TypeString,
                Optional:     true,
                ExactlyOneOf: []string{"name", "id"},
            },
            "id": {
                Type:         schema.TypeString,
                Optional:     true,
                ExactlyOneOf: []string{"name", "id"},
            },
            "description": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "groups": {
                Type:     schema.TypeList,
                Elem:     &schema.Schema{Type: schema.TypeString},
                Computed: true,
            },
        },
    }
}
```

**Features**:
- Lookup by name OR id
- Return description and group memberships
- Use caching with revision tracking (similar to groups datasource)

#### 2.2 Multi-Service User Datasource (Optional)

**File**: `datasources/iam/serviceusers/data_source_multi.go` (new)

Similar to `datasources/iam/groups/data_source_multi.go`

```go
func DataSourceMulti() *schema.Resource {
    return &schema.Resource{
        ReadContext: DataSourceReadMulti,
        Schema: map[string]*schema.Schema{
            "name": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Filter by name (partial match)",
            },
            "service_users": {
                Type:     schema.TypeList,
                Computed: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "id":          { Type: schema.TypeString, Computed: true },
                        "name":        { Type: schema.TypeString, Computed: true },
                        "description": { Type: schema.TypeString, Computed: true },
                    },
                },
            },
        },
    }
}
```

#### 2.3 Provider Registration

**File**: `provider/provider.go`

Add imports and registration:

```go
import (
    ds_iam_serviceusers "github.com/dynatrace-oss/terraform-provider-dynatrace/datasources/iam/serviceusers"
)

// In DataSourcesMap:
"dynatrace_iam_service_user":  ds_iam_serviceusers.DataSource(),
"dynatrace_iam_service_users": ds_iam_serviceusers.DataSourceMulti(), // optional
```

### Phase 3: Export Configuration

**File**: `dynatrace/export/enums.go`

Add enum entries:

```go
// In ResourceTypes struct:
IAMServiceUser ResourceType

// In ResourceType String() method:
case ResourceTypes.IAMServiceUser:
    return "dynatrace_iam_service_user"

// In ValidResources array:
"dynatrace_iam_service_user",
```

### Phase 4: Documentation

#### 5.1 Resource Documentation

**File**: `templates/resources/iam_service_user.md.tmpl` (new)

```
---
page_title: "dynatrace_iam_service_user Resource - terraform-provider-dynatrace"
subcategory: "IAM"
description: |-
  Manages IAM service users for API access and automation
---

# dynatrace_iam_service_user (Resource)

Service users are non-human identities used for programmatic access to Dynatrace.

## Example Usage

```
