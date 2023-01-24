package quota

// HostUnits represents host units consumption and quota information on environment level. If skipped when editing via PUT method then already set quota will remain
type HostUnits struct {
	MaxLimit *int64 `json:"maxLimit"` // Concurrent environment quota. Not set if unlimited. When updating via PUT method, skipping this field will set quota unlimited
}

func (me *HostUnits) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil
}
