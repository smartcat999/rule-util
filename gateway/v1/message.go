package v1

type MsgData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

//type OnvifPushMeta struct {
//	Id       string `json:"id"`
//	Name     string `json:"name"`
//	Address  string `json:"address"`
//}
type OnvifPushMsg struct {
	DriverId string      `json:"driver_id"`
	DeviceId string      `json:"device_id"`
	Version  string      `json:"version"`
	EdgeName string      `json:"edge_name"`
	Meta     interface{} `json:"meta"`
	Time     int64       `json:"time"`
}
