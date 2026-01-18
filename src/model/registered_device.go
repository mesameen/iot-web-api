package model

type RegisteredDevice struct {
	Imei          string `json:"imei"`
	TenantGroupID string `json:"tenant_group_id"`
	ParserID      int32  `json:"parser_id"`
	DeviceTypeID  int32  `json:"device_type_id"`
	TenantID      string `json:"tenant_id"`
	Status        int32  `json:"status"`
}
