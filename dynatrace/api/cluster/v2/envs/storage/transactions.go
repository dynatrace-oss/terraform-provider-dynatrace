package storage

// Transactions represents transaction storage usage and limit information on environment level. If skipped when editing via PUT method then already set limit will remain
type Transactions struct {
	MaxLimit *int64 `json:"maxLimit"` // Maximum storage limit [bytes]
}

func (me *Transactions) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil
}
