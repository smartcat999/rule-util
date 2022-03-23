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

package logf

import (
	metadata "github.com/tkeel-io/rule-util/metadata/v1"
	"go.uber.org/zap"
)

var (
	Any           = zap.Any
	Error         = zap.Error
	String        = zap.String
	Float32       = zap.Float32
	Float64       = zap.Float64
	Float32s      = zap.Float32s
	Float64s      = zap.Float64s
	Int           = zap.Int
	Int32         = zap.Int32
	Int64         = zap.Int64
	Bool          = zap.Bool
	ByteString    = zap.ByteString
	Time          = zap.Time
	Fields        = zap.Fields
	AddStacktrace = zap.AddStacktrace
)

// String constructs a field with the given key and value.
func Message(message metadata.Message) zap.Field {
	return Any("message", message)
}
func UserID(userID string) zap.Field {
	return String("userID", userID)
}
func DeviceID(deviceID string) zap.Field {
	return String("deviceID", deviceID)
}
func EntityID(entityID string) zap.Field {
	return String("entityID", entityID)
}
func ThingID(thingID string) zap.Field {
	return String("thingID", thingID)
}
func PropertyID(propertyID string) zap.Field {
	return String("propertyID", propertyID)
}
func EventID(eventID string) zap.Field {
	return String("eventID", eventID)
}
func ServiceID(serviceID string) zap.Field {
	return String("serviceID", serviceID)
}
func StatusCode(statusCode interface{}) zap.Field {
	return Any("statusCode", statusCode)
}
func Topic(topic string) zap.Field {
	return String("topic", topic)
}
