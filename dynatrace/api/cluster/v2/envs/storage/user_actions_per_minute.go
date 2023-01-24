package storage

// UserActionsPerMinute represents the maximum number of user actions generated per minute on environment level. Can be set to any value from 1 to 2147483646 or left unlimited. If skipped when editing via PUT method then already set limit will remain
type UserActionsPerMinute struct {
	MaxLimit *int32 `json:"maxLimit"` // Maximum traffic [units per minute]
}

func (me *UserActionsPerMinute) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil
}
