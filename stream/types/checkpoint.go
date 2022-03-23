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

package types

import "context"

// ------------------------------------------------------------------------
//  SourceTransport
// ------------------------------------------------------------------------
type SourceTransport interface {
	CheckpointedFunction
	CheckpointListener
	SourceFunction
}

// ------------------------------------------------------------------------
//  SinkAsyncTransport
// ------------------------------------------------------------------------
type SinkAsyncTransport interface {
	CheckpointedFunction
	SinkAsyncFunction
}

// ------------------------------------------------------------------------
//  source function
// ------------------------------------------------------------------------
type SourceFunction interface {
	Open(ctx context.Context, configuration Configuration) error
	Run(ctx context.Context, stx SourceContext) error
	Cancel(ctx context.Context)
	Close(ctx context.Context) error
}

type SourceContext interface {
	Collect(msg interface{})
	MarkAsTemporarilyIdle()
}

// ------------------------------------------------------------------------
//  sink function
// ------------------------------------------------------------------------
type SinkAsyncFunction interface {
	Open(ctx context.Context, configuration Configuration) error
	AsyncInvoke(ctx context.Context, message interface{}, stx SinkContext) error
	Cancel(ctx context.Context)
	Close(ctx context.Context) error
}

type SinkContext interface {
	Promise() Promise
	///** Returns the current processing time. */
	//CurrentProcessingTime() int64
	///** Returns the current event-time watermark. */
	//CurrentWatermark() Promise
	///**
	// * Returns the timestamp of the current input record or {@code null} if the element does not
	// * have an assigned timestamp.
	// */
	//Timestamp() int64
}

// AsyncCollector: For each input streaming record, an AsyncCollector
// will be created and passed into user's callback to get the async i/o result.
type AsyncCollector interface {
}

type AsyncCollectorBuffer interface {
}

// ------------------------------------------------------------------------
//  common
// ------------------------------------------------------------------------
type Configuration interface {
	Name() string
	Metadata() map[string][]byte
}

//type ConfigMap interface {
//	Int(key string) int
//	Int32(key string) int32
//	Int64(key string) int64
//	Float(key string) float32
//	Float32(key string) float32
//	Float64(key string) float64
//	String(key string) string
//	Bool(key string) string
//}

// ------------------------------------------------------------------------
//  check pointed
// ------------------------------------------------------------------------
type CheckpointID int64
type CheckpointTimestamp int64
type CheckpointedFunction interface {
	SnapshotState(context FunctionSnapshotContext) error
	InitializeState(context FunctionInitializationContext) error
}

type CheckpointListener interface {
	NotifyCheckpointComplete(sid CheckpointID) error
}

type FunctionSnapshotContext = ManagedSnapshotContext

type ManagedSnapshotContext interface {
	GetCheckpointId() CheckpointID
	GetCheckpointTimestamp() CheckpointTimestamp
}

//type FunctionInitializationContext = ManagedInitializationContext
type FunctionInitializationContext = ManagedSnapshotContext

type ManagedInitializationContext interface {
	IsRestored() bool
	GetStateStore() StateStore
}

type StateStore interface {
	GetState(stateProperties Configuration) []State
}

type State interface {
	Clear()
}
