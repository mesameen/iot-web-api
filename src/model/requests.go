package model

type GetTelematicsDataRequest struct {
	IMEI string `json:"imei"`
	From int64  `json:"from"`
	To   int64  `json:"to"`
}

type GetConnectionsDataRequest struct {
	IMEI string `json:"imei"`
	From int64  `json:"from"`
	To   int64  `json:"to"`
}

type GetRegisteredDevicesRequest struct {
	IMEI string `json:"imei"`
}
