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

package stream

import (
	gateway "github.com/tkeel-io/rule-util/gateway/v1"
	metadata "github.com/tkeel-io/rule-util/metadata/v1"
	"github.com/tkeel-io/rule-util/stream/checkpoint"
	"github.com/tkeel-io/rule-util/stream/factor"
	"github.com/tkeel-io/rule-util/stream/types"
)

type SubData = gateway.SubData
type UnSubData = gateway.UnSubData

type Message = metadata.Message
type PublishMessage = metadata.PublishMessage
type ConnectMessage = metadata.ConnectMessage
type DisConnectMessage = metadata.DisConnectMessage
type SubMessage = metadata.SubMessage
type UnSubMessage = metadata.UnSubMessage

var NewMessage = metadata.NewMessage

type SinkClient = types.SinkClient
type SourceClient = types.SourceClient
type SourceHandle = types.SourceHandle

type Sink = types.Sink__
type Source = types.Source__

var OpenSink = checkpoint.OpenSink
var OpenSource = checkpoint.OpenSource
var OpenSinkContext = checkpoint.OpenSinkContext
var OpenSourceContext = checkpoint.OpenSourceContext

var ContextKeySendRAW = types.ContextKeySendRAW

var NoopSink = factor.NoopSink
