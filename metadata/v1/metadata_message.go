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

import (
	"fmt"
	_ "github.com/gogo/googleapis/google/api"
	_ "github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	_ "github.com/mwitkow/go-proto-validators"
	"strconv"
	"strings"
	"time"
)

const (
	MATE_VERSION           = "x-stream-version"
	MATE_DOMAIN            = "x-stream-domain"
	MATE_PACKET_IDENTIFIER = "x-stream-packet-id"
	MATE_TYPE              = "x-stream-type"
	MATE_ORIGIN_MESSAGE_ID = "x-stream-origin-id"
	MATE_METHOD            = "x-stream-method"
	MATE_TOPIC             = "x-stream-topic"
	MATE_QOS               = "x-stream-qos"
	MATE_TTL               = "x-stream-ttl"
	MATE_TARGET_ENTITY     = "x-stream-target-entity"
	MATE_TARGET_ENTITIES   = "x-stream-target-entities"
	MATE_ENTITY            = "x-stream-entity"
	MATE_ENTITY_SOURCE     = "x-stream-entity-source"
	//MATE_ENTITY_GROUP      = "x-stream-entity-group"
	MESSAGE_VERSION_1_0    = "1.0"
)

var _ Attribution = &ProtoMessage{}
var _ Message = &ProtoMessage{}
var _ PublishMessage = &ProtoMessage{}
var _ ConnectMessage = &ProtoMessage{}
var _ DisConnectMessage = &ProtoMessage{}
var _ SubMessage = &ProtoMessage{}
var _ UnSubMessage = &ProtoMessage{}

type Timestamp = types.Timestamp

type Attribution interface {
	Attr(key string) []byte
	SetAttr(key string, value []byte) Message
	ForeachAttr(handler func(key, val string) error) error
}

type Message interface {
	Attribution
	Version() string
	Copy() Message
	Validate() error
	Domain() string
	SetDomain(domain string) Message
	Method() string
	SetMethod(method string) Message
	PacketIdentifier() string
	SetPacketIdentifier(id string) Message
	Entity() string
	SetEntity(entityId string) Message

	//Connect
	SetEntitySource(entitySource string) Message
	SetSource(entitySource string) Message

	//Publish
	SetTopic(topic string) Message
	SetQoS(qos int) Message
	// Deprecated: Please use AddTargetEntity
	SetTargetEntity(entityId string) Message
	SetTTL(qos int) Message

	// Action
	//SetTargetDeviceEntity(deviceId string) Message
	AddTargetEntity(entityId string) Message
	SetTargetEntities(entityIds []string) Message
	CleanTargetEntity() Message

	Data() []byte
	SetData(data []byte) Message
}

type PublishMessage interface {
	Message
	Topic() string
	// Deprecated: Please use TargetEntities
	TargetEntity() string
	TargetEntities() []string
	QoS() int
	TTL() int
}

type ConnectMessage interface {
	Message
	Source() string
}
type DisConnectMessage interface {
	Message
	Source() string
}

type SubMessage interface {
	Message
}

type UnSubMessage interface {
	Message
}

type MessageOption func(Message) Message

func NewMessage(opts ...MessageOption) *ProtoMessage {
	now := time.Now()
	return &ProtoMessage{
		Matedata: map[string]string{
			MATE_VERSION: MESSAGE_VERSION_1_0,
		},
		Time: Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	}
}

func (this *ProtoMessage) Copy() Message {
	bt, _ := proto.Marshal(this)
	msg := &ProtoMessage{
	}
	proto.Unmarshal(bt, msg)
	if msg != nil {
		msg.setOrigin(this.PacketIdentifier())
	}
	return msg
}

func (this *ProtoMessage) setOrigin(ID string) Message {
	return this.setMate(MATE_ORIGIN_MESSAGE_ID, ID)
}

func (this *ProtoMessage) setMate(key, value string) Message {
	if this.Matedata == nil {
		this.Matedata = make(map[string]string)
	}
	this.Matedata[key] = value
	return this
}
func (this *ProtoMessage) delMate(key string) Message {
	if this.Matedata == nil {
		this.Matedata = make(map[string]string)
	}
	delete(this.Matedata, key)
	return this
}
func (this *ProtoMessage) getMate(key string) string {
	if this.Matedata == nil {
		this.Matedata = make(map[string]string)
		return ""
	}
	if t, has := this.Matedata[key]; has {
		return t
	}
	return ""
}

func (this *ProtoMessage) Attr(key string) []byte {
	return []byte(this.getMate(key))
}

func (this *ProtoMessage) SetAttr(key string, value []byte) Message {
	if value == nil {
		return this.delMate(key)
	}
	return this.setMate(key, string(value))
}

func (this *ProtoMessage) ForeachAttr(handler func(key, val string) error) error {
	for k, v := range this.Matedata {
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (this *ProtoMessage) Domain() string {
	return this.getMate(MATE_DOMAIN)
}

func (this *ProtoMessage) SetDomain(domain string) Message {
	return this.setMate(MATE_DOMAIN, domain)
}

func (this *ProtoMessage) Method() string {
	return this.getMate(MATE_METHOD)
}

func (this *ProtoMessage) SetMethod(method string) Message {
	return this.setMate(MATE_METHOD, method)
}

func (this *ProtoMessage) PacketIdentifier() string {
	return this.getMate(MATE_PACKET_IDENTIFIER)
}

func (this *ProtoMessage) SetPacketIdentifier(id string) Message {
	return this.setMate(MATE_PACKET_IDENTIFIER, id)
}

func (this *ProtoMessage) Version() string {
	return this.getMate(MATE_VERSION)
}

func (this *ProtoMessage) Topic() string {
	return this.getMate(MATE_TOPIC)
}

func (this *ProtoMessage) SetTopic(topic string) Message {
	return this.setMate(MATE_TOPIC, topic)
}

func (this *ProtoMessage) TargetEntity() string {
	return this.getMate(MATE_TARGET_ENTITY)
}

var EmptyTargetEntities = make([]string, 0)

func (this *ProtoMessage) TargetEntities() []string {
	entities := this.getMate(MATE_TARGET_ENTITIES)
	if entities == "" {
		return EmptyTargetEntities
	} else {
		return strings.Split(entities, ",")
	}
}

func (this *ProtoMessage) SetTargetEntity(entityId string) Message {
	return this.setMate(MATE_TARGET_ENTITY, entityId)
}

func (this *ProtoMessage) SetTargetEntities(entityIds []string) Message {
	i := 0 // output index
	for _, entity := range entityIds {
		if entity != "" {
			entityIds[i] = entity
			i++
		}
	}
	return this.setMate(MATE_TARGET_ENTITIES, strings.Join(entityIds[:i], ","))
}

//TODO remove duplicate entity
func (this *ProtoMessage) AddTargetEntity(entityId string) Message {
	entities := this.getMate(MATE_TARGET_ENTITIES)
	if entities == "" {
		return this.setMate(MATE_TARGET_ENTITIES, fmt.Sprintf("%s", entityId))
	}
	return this.setMate(MATE_TARGET_ENTITIES, fmt.Sprintf("%s,%s", entities, entityId))
}

func (this *ProtoMessage) CleanTargetEntity() Message {
	return this.delMate(MATE_TARGET_ENTITIES)
}

func (this *ProtoMessage) QoS() int {
	qos := this.getMate(MATE_QOS)
	if i, err := strconv.ParseInt(qos, 10, 32); err == nil {
		return int(i)
	}
	return 0
}

func (this *ProtoMessage) SetQoS(qos int) Message {
	return this.setMate(MATE_QOS, fmt.Sprintf("%d", qos))
}

func (this *ProtoMessage) TTL() int {
	qos := this.getMate(MATE_TTL)
	if qos == "" {
		return 0
	}
	if i, err := strconv.ParseInt(qos, 10, 32); err == nil {
		return int(i)
	}
	return 0
}

func (this *ProtoMessage) SetTTL(ttl int) Message {
	return this.setMate(MATE_TTL, fmt.Sprintf("%d", ttl))
}

func (this *ProtoMessage) Data() []byte {
	return this.RawData
}

func (this *ProtoMessage) SetData(data []byte) Message {
	this.RawData = data
	return this
}

func (this *ProtoMessage) SetEntity(entityId string) Message {
	return this.setMate(MATE_ENTITY, entityId)
}

func (this *ProtoMessage) Entity() string {
	return this.getMate(MATE_ENTITY)
}

func (this *ProtoMessage) SetSource(entity string) Message {
	return this.setMate(MATE_ENTITY_SOURCE, entity)
}

func (this *ProtoMessage) SetEntitySource(entity string) Message {
	return this.setMate(MATE_ENTITY_SOURCE, entity)
}
func (this *ProtoMessage) Source() string {
	return this.getMate(MATE_ENTITY_SOURCE)
}


//func (this *ProtoMessage) SetEntityGroup(group string) Message {
//	return this.setMate(MATE_ENTITY_GROUP, group)
//}
//
//func (this *ProtoMessage) EntityGroup() string {
//	return this.getMate(MATE_ENTITY_GROUP)
//}
