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

package kafka

import (
	"github.com/Shopify/sarama"
	"time"
)

type ConsumerOption struct {
	Addrs    string `json:"addrs,omitempty"`
	ClientID string `json:"client_id,omitempty"`
	GroupID  string `json:"groupID,omitempty"`
	Topics   string `json:"topics,omitempty"`
}

const SinkSendRetryCnt = 3
const SinkSendRetryInterval = 1000 * time.Millisecond
const MaxReconnectTime = 5 * time.Minute

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}

//type Status int
//const (
//	IDEL = Status(iota)
//	RUNNING
//)

type SourceOption struct {
	Addrs    []string `json:"addrs,omitempty"`
	MsgMode  string   `json:"msgMode,omitempty"`
	LogLevel string   `json:"logLevel,omitempty"`
	Topic    string   `json:"topic,omitempty"`
	GroupID  string   `json:"groupID,omitempty"`
}

type SinkOption struct {
	Addrs         []string `json:"addrs,omitempty"`
	MsgMode       string   `json:"msgMode,omitempty"`
	LogLevel      string   `json:"logLevel,omitempty"`
	Topic         string   `json:"topic,omitempty"`
	GroupID       string   `json:"groupID,omitempty"`
	retryCnt      int
	retryInterval time.Duration
	config        *sarama.Config
}

type ackContext struct {
	offset int64
}

type snapshot map[int32]*ackContext
