package retention

// SessionReplay represents session replay retention settings on environment level. Can be set to any value from 1 to 35 days. If skipped when editing via PUT method then already set limit will remain
type SessionReplay struct {
	MaxLimitInDays int64 `json:"maxLimitInDays"` // Maximum retention limit [days]
}
