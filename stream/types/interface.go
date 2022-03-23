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

import (
	"context"
	"net/url"
)

type SourceHandle func(ctx context.Context, message interface{}) error

// Sink
type Sink__ interface {
	String() string
	Send(ctx context.Context, event interface{}) error
	SendAsync(ctx context.Context, event interface{}) (Promise, error)
	Close(ctx context.Context) error
}

// Source
type Source__ interface {
	String() string
	StartReceiver(ctx context.Context, fn SourceHandle) error
	Close(ctx context.Context) error
}

// Client interface defines the runtime contract the Message client supports.
type SourceClient interface {
	StartReceiver(ctx context.Context, fn SourceHandle) error
	Close(ctx context.Context) error
}

// Client interface defines the runtime contract the Message client supports.
type SinkClient interface {
	SendAsync(ctx context.Context, event interface{}) (Promise, error)
	Close(ctx context.Context) error
}

type Promise interface {
	Then(s func(err error)) Promise
	Finish(err error) Promise
}
type Receiver interface {
	Receive(context.Context, interface{}) error
}

//type SourceTransport interface {
//	SetReceiver(Receiver)
//	StartReceiver(ctx context.Context) error
//	Close(ctx context.Context) error
//}
//
//type SinkAsyncTransport interface {
//	Send(ctx context.Context, event Message) (error)
//	SendAsync(ctx context.Context, event Message) Promise
//	Close(ctx context.Context) error
//}

type Driver interface {
	NewSourceTransport(ctx context.Context, url *url.URL) (SourceTransport, error)
	NewSinkTransport(ctx context.Context, url *url.URL) (SinkAsyncTransport, error)
}

// Codec is the interface for transport codecs to convert between transport
// specific payloads and the Message interface.
type Codec interface {
	Encode(interface{}) ([]byte, error)
}
