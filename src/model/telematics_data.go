package model

type TelematicsData struct {
	Imei             string       `json:"imei"`
	GpsData          *GpsData     `json:"gps_data"`
	EventType        string       `json:"event_type"`
	SensorData       *SensorData  `json:"sensor_data"`
	NetworkData      *NetworkData `json:"network_data"`
	DeviceDatetime   uint64       `json:"device_datetime"`
	InsertDatetime   int          `json:"insert_datetime"`
	TenantGroupID    string       `json:"tenant_group_id"`
	ListenerDatetime uint64       `json:"listener_datetime"`
}

type GpsData struct {
	Speed     float32 `json:"speed"`
	Heading   float32 `json:"heading"`
	Altitude  float32 `json:"altitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SensorData struct {
	Rpm          float32 `json:"rpm"`
	Speed        float32 `json:"speed"`
	Idling       bool    `json:"idling"`
	Ibutton      string  `json:"ibutton"`
	Distance     float32 `json:"distance"`
	EngTemp      float32 `json:"eng_temp"`
	Ignition     bool    `json:"ignition"`
	BatteryPer   float32 `json:"battery_per"`
	AccPedalPer  float32 `json:"acc_pedal_per"`
	FuelLevelPer float32 `json:"fuel_level_per"`
}

type NetworkData struct {
	SignalPer   float32 `json:"signal_per"`
	NetworkType string  `json:"network_type"`
}
