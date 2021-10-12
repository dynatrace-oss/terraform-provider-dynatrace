package locations

// 	SyntheticLocations is a list of synthetic locations
type SyntheticLocations struct {
	Locations []LocationCollectionElement `json:"locations"` // A list of synthetic locations
}

type LocationCollectionElement struct {
	ID            string         `json:"entityId"`      // The Dynatrace entity ID of the location
	Name          string         `json:"name"`          // The name of the location
	Type          LocationType   `json:"type"`          // The type of the location
	Status        *Status        `json:"status"`        // The status of the location: \n\n* `ENABLED`: The location is displayed as active in the UI. You can assign monitors to the location. \n* `DISABLED`: The location is displayed as inactive in the UI. You can't assign monitors to the location. Monitors already assigned to the location will stay there and will be executed from the location. \n* `HIDDEN`: The location is not displayed in the UI. You can't assign monitors to the location. You can only set location as `HIDDEN` when no monitor is assigned to it
	CloudPlatform *CloudPlatform `json:"cloudPlatform"` // The cloud provider where the location is hosted. \n\n Only applicable to `PUBLIC` locations
	IPs           []string       `json:"ips"`           // The list of IP addresses assigned to the location. \n\n Only applicable to `PUBLIC` locations
	Stage         *Stage         `json:"stage"`         // The release stage of the location

}
