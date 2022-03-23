/*
 * Copyright (C) 2019 Yunify, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this work except in compliance with the License.
 * You may obtain a copy of the License in the LICENSE file, or at:
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1

import "time"

const (
	THING_PROPERTY_PLATFORM_URL       = "/sys/%s/%s/thing/property/platform/post"
	THING_PROPERTY_PLATFORM_REPLY_URL = "/sys/%s/%s/thing/property/platform/post_reply"
	THING_PROPERTY_BASE_URL           = "/sys/%s/%s/thing/property/base/post"
	THING_PROPERTY_BASE_REPLY_URL     = "/sys/%s/%s/thing/property/base/post_reply"
	THING_EVENT_URL                   = "/sys/%s/%s/thing/event/%s/post"
	THING_EVENT_REPLY_URL             = "/sys/%s/%s/thing/event/%s/post_reply"
	THING_SERVICE_URL                 = "/sys/%s/%s/thing/service/%s/call"
	THING_SERVICE_URL_REPLY           = "/sys/%s/%s/thing/service/%s/call_reply"
)
const (
	USER_UNKNOW         = 0
	USER_MESSAGE_REPORT = 1
	USER_MESSAGE_REPLY  = 2
	USER_MESSAGE_DOWN   = 3
)

type ThingPropertyMeta struct {
	Tags   map[string]string `json:"tags"`
	UserId string            `json:"user_id"`
	Topic  string            `json:"topic"`
	Data   []byte            `json:"data"`
}

//
//up property
type ThingPropertyMsg struct {
	Id       string                   `json:"id"`
	Version  string                   `json:"version"`
	Metadata *DownMeta                `json:"metadata"`
	Type     string                   `json:"type"`
	Params   map[string]*PropertyData `json:"params"`
}

//
////up event
type ThingEventMsg struct {
	Id       string     `json:"id"`
	Version  string     `json:"version"`
	Metadata *DownMeta  `json:"metadata"`
	Type     string     `json:"type"`
	Params   *EventData `json:"params"`
}

type PropertyData struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}

//
type EventData struct {
	Value map[string]interface{} `json:"value"`
	Time  int64                  `json:"time"`
}

// 设备上下线消息
type ThingStatusMsg struct {
	Status       string `json:"status"`
	ThingID      string `json:"thing_id"`
	DeviceID     string `json:"device_id"`
	Time         int64  `json:"time"`
	ConnectionID uint64 `json:"connection_id"`
	NodeName     string `json:"node_name"`
	RemoteAddr   string `json:"remote_addr"`
	InstanceID   string `json:"instance_id"`
}
type DownMeta struct {
	EntityId  string                 `json:"entityId"`
	ModelId   string                 `json:"modelId"`
	SourceId  []string               `json:"sourceId,omitempty"`
	EpochTime int64                  `json:"epochTime,omitempty"`
	Tags      map[string]interface{} `json:"tags,omitempty"`
}

//down
type DownThingMsg struct {
	Id       string                 `json:"id"`
	Version  string                 `json:"version"`
	Metadata *DownMeta              `json:"metadata"`
	Type     string                 `json:"type"`
	Params   map[string]interface{} `json:"params"`
}
type DeviceNotifyMsg struct {
	Id       string      `json:"id"`
	Version  string      `json:"version"`
	Metadata *DownMeta   `json:"metadata"`
	Type     string      `json:"type"`
	Params   interface{} `json:"params"`
}

//down
type GetThingMsg struct {
	Id       string    `json:"id"`
	Version  string    `json:"version"`
	Metadata *DownMeta `json:"metadata"`
	Type     string    `json:"type"`
	Params   []string  `json:"params"`
}
type Reply struct {
	Code int         `json:"code"`
	Id   string      `json:"id"`
	Data interface{} `json:"data"`
}
type ThingProperty struct {
	Id         string                 `json:"id"`
	Name       string                 `json:"name"`
	Identifier string                 `json:"identifier"` //识别字符
	Type       string                 `json:"type"`       //数据类型["int32","float","double""array","bool","enum","date","struct","string"]
	UserId     string                 `json:"user_id"`
	Define     map[string]interface{} `json:"define"`
}
type ThingEvent struct {
	Id         string           `json:"id"`
	Identifier string           `json:"identifier"` //识别字符
	Type       string           `json:"type"`       //数据类型["info","","warn""error""]
	Output     []*ThingProperty `json:"output"`
	UserId     string           `json:"user_id"`
}

//create
type DeviceModel struct {
	UserId           string    `db:"user_id"`
	ThingId          string    `db:"thing_id"`
	DeviceId         string    `db:"device_id"`
	SourceId         string    `db:"source_id"`
	Identifier       string    `db:"identifier"`
	ValueInt32       int32     `db:"value_int32"`
	ValueFloat       float32   `db:"value_float"`
	ValueDouble      float64   `db:"value_double"`
	ValueEnum        int8      `db:"value_enum"`
	ValueString      string    `db:"value_string"`
	ValueBool        bool      `db:"value_bool"`
	ValueInt64       int64     `db:"value_int64"`
	ValueStringEx    string    `db:"value_string_ex"`
	ValueArrayString []string  `db:"value_array_string"`
	ValueArrayStruct []string  `db:"value_array_struct"`
	ValueArrayInt32  []int32   `db:"value_array_int"`
	ValueArrayFloat  []float32 `db:"value_array_float"`
	ValueArrayDouble []float64 `db:"value_array_double"`
	Time             int64     `db:"time"`
	Tags             string    `db:"tags"`
	ActionDate       time.Time `db:"action_date"`
	ActionTime       time.Time `db:"action_time"`
}

//user message
type UserDeviceModel struct {
	MsgId        string    `db:"msg_id"`
	UserId       string    `db:"user_id"`       //用户id
	SourceId     string    `db:"source_id"`     //平台定义设备id
	SerialNumber string    `db:"serial_number"` //设备序列号
	Topic        string    `db:"topic"`         //topic
	Value        []byte    `db:"value"`         //设备消息数据
	Type         uint8     `db:"type"`
	Tags         string    `db:"tags"`
	Time         int64     `db:"time"` //消息时间
	ActionDate   time.Time `db:"action_date"`
	ActionTime   time.Time `db:"action_time"`
}
type EventThingStore struct {
	DeviceId   string    `json:"device_id"`
	ThingId    string    `json:"thing_id"`
	UserId     string    `json:"user_id"`
	SourceId   string    `json:"source_id"`
	Identifier string    `json:"identifier"`
	EventId    string    `json:"event_id"`
	Type       string    `json:"type"`
	Metadata   string    `json:"metadata"`
	Time       int64     `json:"time"`
	ActionDate time.Time `json:"action_date"`
	ActionTime time.Time `json:"action_time"`
}

//设备上下线记录
type StatusHistory struct {
	DeviceId     string    `db:"device_id"`
	ThingId      string    `db:"thing_id"`
	UserId       string    `db:"user_id"` //用户id
	Status       string    `db:"status"`
	ConnectionId int64     `db:"connection_id"` //连接id
	NodeName     string    `db:"node_name"`     //节点id
	RemoteAddr   string    `db:"remote_addr"`   //远程地址
	Time         int64     `db:"time"`          //消息时间
	ActionDate   time.Time `db:"action_date"`
	ActionTime   time.Time `db:"action_time"`
}

type UserMessage struct {
	DeviceId string `db:"device_id"`
	Time     int64  `db:"time"`
	Data     []byte `db:"data"`
	Type     uint8  `db:"type"`
}

type DataProperty struct {
	DeviceId string            `db:"device_id"`
	Time     int64             `db:"time"`
	Data     map[string]string `db:"data"`
}
