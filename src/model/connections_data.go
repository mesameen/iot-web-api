package model

type ConnectionsData struct {
	Imei           string `json:"imei"`
	TenantGroupID  string `json:"tenant_group_id"`
	TenantID       string `json:"tenant_id"`
	ConnectedAt    int64  `json:"connected_at_ms"`
	DisconnectedAt int64  `json:"disconnected_at_ms"`
	Duration       int32  `json:"duration"`
	Addr           string `json:"addr"`
	ListenerName   string `json:"listener_name"`
	Reason         string `json:"reason"`
	Sent           int64  `json:"sent"`
	Recv           int64  `json:"recv"`
	Action         string `json:"action"`
}
