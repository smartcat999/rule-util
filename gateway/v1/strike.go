package v1

type StrikeStatus struct {
	SessionCount  int    `json:"session_count"`
	CollectTime   int64  `json:"collect_time"`
	LastErrorTime int64  `json:"last_error_time"`
	Load          string `json:"load"`
	NodeName      string `json:"node_name"`
	InstanceId    string `json:"instance_id"`
	//	NodeIp        string `json:"node_ip"`
	Uptime int64 `json:"uptime"`
}
