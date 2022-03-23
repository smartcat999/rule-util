/**
 * @Author: hexing
 * @Description:
 * @File:  edge-monitor
 * @Version: 1.0.0
 * @Date: 20-5-26 上午11:40
 */

package v1

const (
	CommitLogProperty_AppId  = "appID"
	CommitLogProperty_Status = "statusCode"
	CommitLogProperty_Source = "source"
)

//log source
const (
	CommitLogSource_EdgeWize = "EdgeWize"
	CommitLogSource_Monitor  = "EdgeMonitor"
)

const (
	//Monitor服务启动
	LogStatus_MonitorStart = 920
	//Monitor服务停止
	LogStatus_MonitorStop = 921
	//Monitor服务运行失败
	LogStatus_MonitorFailed = 922
	//采集数据成功
	LogStatus_ColllectSuccess = 923
	//采集数据失败
	LogStatus_CollectFailed = 924
	//上报数据成功
	LogStatus_PostSuccess = 925
	//上报数据失败
	LogStatus_PostFailed = 926
)

const (
	//DEVICE_STATUS
	DEVICE_STATUS_END  = "iote-global-onoffline-end"
	DEVICE_STATUS_EDGE = "iote-global-onoffline-edge"
	//EXTENSION
	EXTENSION_USERID     = "user_id"
	EXTENSION_THINGID    = "thing_id"
	EXTENSION_DEVICEID   = "device_id"
	EXTENSION_PROPERTYID = "property_id"
)
const (
	PROPERTY_TYPE_POST = "thing.property.post"
	THING_VERSION      = "v1.0"
)
const (
	CPU_PERCENT     = "cpuPercent"
	MEM_USAGE       = "memUsage"
	DISK_USAGE      = "diskUsage"
	BYTES_RECV      = "netIn"
	BYTES_SEND      = "netOut"
	NET_INFO        = "netInfo"
	UP_RATE         = "upRate"
	DOWN_RATE       = "downRate"
	SYSTEM_VERSION  = "system"
	EDGE_VERSION    = "edgeVersion"
	SYSTEM_HARDWARE = "hardware"
)
