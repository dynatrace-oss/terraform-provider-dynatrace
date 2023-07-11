package processgroupalerting

type AlertingMode string

var AlertingModes = struct {
	OnPgiUnavailability      AlertingMode
	OnInstanceCountViolation AlertingMode
}{
	AlertingMode("ON_PGI_UNAVAILABILITY"),
	AlertingMode("ON_INSTANCE_COUNT_VIOLATION"),
}
