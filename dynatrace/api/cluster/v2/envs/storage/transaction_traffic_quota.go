package storage

// TransactionTrafficQuota represetnts the maximum number of newly monitored entry point PurePaths captured per process/minute on environment level. Can be set to any value from 100 to 100000. If skipped when editing via PUT method then already set limit will remain
type TransactionTrafficQuota struct {
	MaxLimit *int32 `json:"maxLimit"` // Maximum traffic [units per minute]
}

func (me *TransactionTrafficQuota) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil
}
