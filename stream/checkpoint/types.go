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

package checkpoint

import (
	"github.com/tkeel-io/rule-util/stream/types"

	_ "github.com/tkeel-io/rule-util/stream/transport/kafka"
)

// Use Cloudevents
//var EventCodec = CodecCloudeventsPbMessage
//// Remove Cloudevents
var EventCodec = CodecPbMessage

type Attribution = types.Attribution
type Message = types.Message
type PublishMessage = types.PublishMessage
type Promise = types.Promise
type SourceHandle = types.SourceHandle

var NewMessage = types.NewMessage

type Sink = types.Sink__
type Source = types.Source__
