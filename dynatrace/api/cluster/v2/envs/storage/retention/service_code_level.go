package retention

// ServiceCodeLevel represents service code level retention settings on environment level. Service code level retention time can't be greater than service request level retention time and both can't exceed one year.If skipped when editing via PUT method then already set limit will remain
type ServiceCodeLevel struct {
	MaxLimitInDays int64 `json:"maxLimitInDays"` // Maximum retention limit [days]
}
