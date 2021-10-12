package locations

import "github.com/dtcookie/hcl"

type Locations []*Location

func (me Locations) Schema() map[string]*hcl.Schema {
	return map[string]*hcl.Schema{
		"location": {
			Type:        hcl.TypeList,
			Description: "The name of the location",
			Optional:    true,
			Elem:        &hcl.Resource{Schema: new(Location).Schema()},
		},
	}
}

func (me *Locations) MarshalHCL() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	entries := []interface{}{}
	for _, entry := range *me {
		if marshalled, err := entry.MarshalHCL(); err == nil {
			entries = append(entries, marshalled)
		} else {
			return nil, err
		}
	}
	result["location"] = entries
	return result, nil
}

type Location struct {
	ID            string         `json:"id"`            // The unique id of the location
	Name          string         `json:"name"`          // The name of the location
	Type          LocationType   `json:"type"`          // Defines the actual set of fields depending on the value. See one of the following objects: \n\n* `PUBLIC` -> PublicSyntheticLocation \n* `PRIVATE` -> PrivateSyntheticLocation \n* `CLUSTER` -> PrivateSyntheticLocation
	Status        *Status        `json:"status"`        // The status of the location: \n\n* `ENABLED`: The location is displayed as active in the UI. You can assign monitors to the location. \n* `DISABLED`: The location is displayed as inactive in the UI. You can't assign monitors to the location. Monitors already assigned to the location will stay there and will be executed from the location. \n* `HIDDEN`: The location is not displayed in the UI. You can't assign monitors to the location. You can only set location as `HIDDEN` when no monitor is assigned to it
	CloudPlatform *CloudPlatform `json:"cloudPlatform"` // The cloud provider where the location is hosted. \n\n Only applicable to `PUBLIC` locations
	IPs           []string       `json:"ips"`           // The list of IP addresses assigned to the location. \n\n Only applicable to `PUBLIC` locations
	Stage         *Stage         `json:"stage"`         // The release stage of the location
}

func (me *Location) Schema() map[string]*hcl.Schema {
	return map[string]*hcl.Schema{
		"id": {
			Type:        hcl.TypeString,
			Description: "The unique ID of the location",
			Optional:    true,
		},
		"name": {
			Type:        hcl.TypeString,
			Description: "The name of the location",
			Optional:    true,
		},
		"type": {
			Type:        hcl.TypeString,
			Description: "The type of the location. Supported values are `PUBLIC`, `PRIVATE` and `CLUSTER`",
			Optional:    true,
		},
		"status": {
			Type:        hcl.TypeString,
			Description: "The status of the location: \n\n* `ENABLED`: The location is displayed as active in the UI. You can assign monitors to the location. \n* `DISABLED`: The location is displayed as inactive in the UI. You can't assign monitors to the location. Monitors already assigned to the location will stay there and will be executed from the location. \n* `HIDDEN`: The location is not displayed in the UI. You can't assign monitors to the location. You can only set location as `HIDDEN` when no monitor is assigned to it",
			Optional:    true,
			Computed:    true,
		},
		"stage": {
			Type:        hcl.TypeString,
			Description: "The release stage of the location",
			Optional:    true,
			Computed:    true,
		},
		"cloud_platform": {
			Type:        hcl.TypeString,
			Description: "The cloud provider where the location is hosted. \n\n Only applicable to `PUBLIC` locations",
			Optional:    true,
			Computed:    true,
		},
		"ips": {
			Type:        hcl.TypeList,
			Description: "The list of IP addresses assigned to the location. \n\n Only applicable to `PUBLIC` locations",
			Elem:        &hcl.Schema{Type: hcl.TypeString},
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Location) MarshalHCL() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	result["id"] = me.ID
	result["name"] = me.Name
	result["type"] = string(me.Type)
	if me.Stage != nil {
		result["stage"] = string(*me.Stage)
	}
	if me.Status != nil {
		result["status"] = string(*me.Status)
	}
	if me.CloudPlatform != nil {
		result["cloud_platform"] = string(*me.CloudPlatform)
	}
	if me.IPs != nil {
		result["ips"] = append([]string{}, me.IPs...)
	}
	return result, nil
}

type Stage string

var Stages = struct {
	Beta       Stage
	ComingSoon Stage
	GA         Stage
}{
	Stage("BETA"),
	Stage("COMING_SOON"),
	Stage("GA"),
}

type CloudPlatform string

var CloudPlatforms = struct {
	Alibaba        CloudPlatform
	AmaconEC2      CloudPlatform
	Azure          CloudPlatform
	DynatraceCloud CloudPlatform
	GoogleCloud    CloudPlatform
	Interoute      CloudPlatform
	Other          CloudPlatform
	Undefined      CloudPlatform
}{
	CloudPlatform("ALIBABA"),
	CloudPlatform("AMAZON_EC2"),
	CloudPlatform("AZURE"),
	CloudPlatform("DYNATRACE_CLOUD"),
	CloudPlatform("GOOGLE_CLOUD"),
	CloudPlatform("INTEROUTE"),
	CloudPlatform("OTHER"),
	CloudPlatform("UNDEFINED"),
}

type Status string

var Statuses = struct {
	Disabled Status
	Enabled  Status
	Hidden   Status
}{
	Status("DISABLED"),
	Status("ENABLED"),
	Status("HIDDEN"),
}

type LocationType string

var LocationTypes = struct {
	Public  LocationType
	Private LocationType
	Cluster LocationType
}{
	LocationType("PUBLIC"),
	LocationType("PRIVATE"),
	LocationType("CLUSTER"),
}
