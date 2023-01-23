package retention

// LogMonitoring represents log monitoring retention settings on environment level. Not editable when Log monitoring is not allowed by license or not configured on cluster level. Can be set to any value from 5 to 90 days. If skipped when editing via PUT method then already set limit will remain
type LogMonitoring struct {
	MaxLimitInDays int64 `json:"maxLimitInDays"` // Maximum retention limit [days]
}
