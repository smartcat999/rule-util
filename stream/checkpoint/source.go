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
	"context"
	"net/url"
	"sync"

	"github.com/tkeel-io/rule-util/pkg/log"
	logf "github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/errors"
	"github.com/tkeel-io/rule-util/stream/factor"
	"github.com/tkeel-io/rule-util/stream/transport"
	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/opentracing/opentracing-go"
)

var _ types.Source__ = (*sourceOperation)(nil)

type sourceOperation struct {
	source            types.SourceTransport
	startOnce         sync.Once
	callbacks         []SourceHandle
	lock              sync.Mutex
	checkpointContext types.ManagedSnapshotContext
	running           bool
	dsn               string
}

func OpenSource(dataSourceName string) (Source, error) {
	return OpenSourceContext(context.Background(), dataSourceName)
}

func OpenSourceContext(ctx context.Context, dataSourceName string) (Source, error) {
	source := factor.NewLoopbackFromDNS(ctx, dataSourceName)
	if source != nil {
		return source, nil
	}
	if u, err := url.Parse(dataSourceName); err != nil {
		return nil, errors.Errorf("stream: failed to parse data source name %q ", dataSourceName)
	} else if drv, ok := transport.GetDriver(u.Scheme); !ok {
		return nil, errors.Errorf("stream: unknown driver %q %s ", u.Scheme, u.String())
	} else if st, err := drv.NewSourceTransport(ctx, u); err != nil {
		return nil, errors.Errorf("stream: failed to create source %q %s ", u.Scheme, u.String())
	} else {
		sourceOP := &sourceOperation{source: st, running: false, dsn: dataSourceName}
		AddSource(st)
		return sourceOP, nil
	}
}

func (this *sourceOperation) String() string {
	return this.dsn
}

func (this *sourceOperation) StartReceiver(ctx context.Context, fn SourceHandle) error {
	this.lock.Lock()
	this.callbacks = append(this.callbacks, fn)
	if !this.running {
		go func() {
			if err := this.source.Run(ctx, this); err != nil {
				log.Error("source Run",
					logf.Error(err))
			}
		}()
		this.running = true
	}
	this.lock.Unlock()
	return nil
}

func (this *sourceOperation) Close(ctx context.Context) error {
	return this.source.Close(ctx)
}

//Collect  SourceContext/Collect
func (this *sourceOperation) Collect(msg interface{}) {
	ctx := context.Background()
	//event := new(PacketMessage)
	if byt, ok := msg.([]byte); ok {
		message, err := EventCodec.Decode(byt)
		if err == nil && len(message.Method()) != 0 {
			span := SpanFromEvent(message, "Receive")
			ctx = ContextWithSpan(ctx, span)
			defer span.Finish()
			this.lock.Lock()
			defer this.lock.Unlock()
			this.put(ctx, message)
			return
		}
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.put(ctx, msg)
	return
}

func ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	return opentracing.ContextWithSpan(ctx, span)
}

//MarkAsTemporarilyIdle  SourceContext/MarkAsTemporarilyIdle
func (this *sourceOperation) MarkAsTemporarilyIdle() {
	panic("implement me")
}

func (this *sourceOperation) put(ctx context.Context, message interface{}) {
	for _, cb := range this.callbacks {
		if err := cb(ctx, message); err != nil {
			log.Error("sourceOperation", logf.Error(err))
		}
	}
}

func (this *sourceOperation) InitializeState(ctx types.FunctionInitializationContext) error {
	return this.source.InitializeState(ctx)
}

func (this *sourceOperation) SnapshotState(checkpoint types.ManagedSnapshotContext) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.source.SnapshotState(checkpoint)
}
