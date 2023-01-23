package storage

// LogMonitoringStorage represents log monitoring storage usage and limit information on environment level. Not editable when Log monitoring is not allowed by license or not configured on cluster level. If skipped when editing via PUT method then already set limit will remain
type LogMonitoringStorage struct {
	MaxLimit *int64 `json:"maxLimit"` // Maximum storage limit [bytes]
}

func (me *LogMonitoringStorage) IsEmpty() bool {
	return me == nil || me.MaxLimit == nil
}
