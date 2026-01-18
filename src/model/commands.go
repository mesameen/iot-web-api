package model

type Command struct {
	ID                 int64  `json:"id"`
	Imei               string `json:"imei"`
	Data               string `json:"data"`
	IsResponseRequired bool   `json:"is_response_required"`
	Response           string `json:"response"`
	MaxRetries         int32  `json:"max_retries"`
	ExpiresAtMs        int64  `json:"expires_at_ms"`
	RetriesCount       int32  `json:"retries_count"`
	TenantGroupID      string `json:"tenant_group_id"`
	SentToDevice       bool   `json:"sent_to_device"`
	SentAtMs           int64  `json:"sent_at_ms"`
	ResponseAtMs       int64  `json:"response_at_ms"`
	TenantID           string `json:"tenant_id"`
}
