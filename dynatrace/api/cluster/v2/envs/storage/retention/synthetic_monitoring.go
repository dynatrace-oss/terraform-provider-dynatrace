package retention

// SyntheticMonitoring represents synthetic monitoring retention settings on environment level. Can be set to any value from 1 to 35 days. If skipped when editing via PUT method then already set limit will remain
type SyntheticMonitoring struct {
	MaxLimitInDays int64 `json:"maxLimitInDays"` // Maximum retention limit [days]
}
