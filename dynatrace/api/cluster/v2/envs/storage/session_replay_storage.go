package storage

// SessionReplayStorage represents session replay storage usage and limit information on environment level. If skipped when editing via PUT method then already set limit will remain
type SessionReplayStorage struct {
	MaxLimit *int64 `json:"maxLimit"` // Maximum storage limit [bytes]
}

func (me *SessionReplayStorage) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil
}
