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

	"github.com/tkeel-io/rule-util/pkg/log"
	logf "github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/errors"
	"github.com/tkeel-io/rule-util/stream/factor"
	"github.com/tkeel-io/rule-util/stream/transport"
	"github.com/tkeel-io/rule-util/stream/types"
)

var _ types.Sink__ = (*sinkOperation)(nil)

type sinkOperation struct {
	sink types.SinkAsyncTransport
	dsn  string
}

func OpenSink(dataSourceName string) (Sink, error) {
	return OpenSinkContext(context.Background(), dataSourceName)
}

func OpenSinkContext(ctx context.Context, dataSourceName string) (Sink, error) {
	if sink := factor.NewLoopbackFromDNS(ctx, dataSourceName); sink != nil {
		return sink, nil
	}

	if sink := factor.NewNoopStreamFromDNS(ctx, dataSourceName); sink != nil {
		return sink, nil
	}

	if u, err := url.Parse(dataSourceName); err != nil {
		return nil, errors.Errorf("stream: failed to parse data source name %q ", dataSourceName)
	} else if drv, ok := transport.GetDriver(u.Scheme); !ok {
		return nil, errors.Errorf("stream: unknown driver %q %s ", u.Scheme, u.String())
	} else if st, err := drv.NewSinkTransport(ctx, u); err != nil {
		return nil, errors.Errorf("stream: failed to create sink %q %s ", u.Scheme, u.String())
	} else {
		AddSink(st)
		return &sinkOperation{
			dsn:  dataSourceName,
			sink: st,
		}, nil
	}
}

func (this *sinkOperation) String() string {
	return this.dsn
}

func (this *sinkOperation) Send(ctx context.Context, message interface{}) error {
	stx := types.NewSinkContext()
	stx.Promise().Then(func(err error) {
		if err != nil {
			log.Error("Send fail",
				logf.Any("message", message),
				logf.Error(err))
		}
	})

	span, ctx := StartSpanFromContext(ctx, "Send")
	defer span.Finish()
	EventWithSpan(message, span)

	return this.sink.AsyncInvoke(ctx, message, stx)
}

func (sink *sinkOperation) SendAsync(ctx context.Context, message interface{}) (Promise, error) {
	stx := types.NewSinkContext()
	stx.Promise().Then(func(err error) {
		log.Error("Send fail",
			logf.Any("message", message),
			logf.Error(err))
	})

	span, ctx := StartSpanFromContext(ctx, "Send")
	defer span.Finish()
	EventWithSpan(message, span)

	return stx.Promise(), sink.sink.AsyncInvoke(ctx, message, stx)
}

func (sink *sinkOperation) Close(ctx context.Context) error {
	return sink.sink.Close(ctx)
}
