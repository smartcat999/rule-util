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

/*
 *  dataSourceName:
 *  qmq://[user[:password]@][net[(addr)]]/{topic}/{group_id}[?param1=value1&paramN=valueN]
 *
 */

package client

import (
	"context"
	"github.com/tkeel-io/rule-util/stream/errors"
	"github.com/tkeel-io/rule-util/stream/transport"
	"github.com/tkeel-io/rule-util/stream/types"
	"net/url"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]types.Driver)
)

// OpenSource open cloudevents channel to receiver msg.
func OpenSource(dataSourceName string, opts ...Option) (source types.SourceClient, err error) {
	return OpenSourceContext(context.Background(), dataSourceName, opts...)
}

// OpenSource open cloudevents channel to receiver msg.
func OpenSourceContext(ctx context.Context, dataSourceName string, opts ...Option) (source types.SourceClient, err error) {
	var (
		u *url.URL
		t types.SourceTransport
	)
	u, err = url.Parse(dataSourceName)
	if err != nil {
		return nil, errors.Errorf("cloudevent: failed to parse data source name %q ", dataSourceName)
	}
	driveri, ok := transport.GetDriver(u.Scheme)
	if !ok {
		return nil, errors.Errorf("cloudevent: unknown driver %q %s ", u.Scheme, u.String())
	}

	t, err = driveri.NewSourceTransport(ctx, u)
	if err != nil {
		return nil, errors.Errorf("cloudevent: failed to create source %q %s ", u.Scheme, u.String())
	}

	return NewSource(ctx, t, opts...)
}


// OpenSource open cloudevents channel to receiver msg.
func OpenSink(dataSourceName string, opts ...Option) (source types.SinkClient, err error) {
	return OpenSinkContext(context.Background(), dataSourceName, opts...)
}
// OpenSink open cloudevents channel to publish msg.
func OpenSinkContext(ctx context.Context, dataSourceName string, opts ...Option) (sink types.SinkClient, err error) {
	var (
		u *url.URL
		t types.SinkAsyncTransport
	)
	u, err = url.Parse(dataSourceName)
	if err != nil {
		return nil, errors.Errorf("cloudevent: failed to parse data source name %q ", dataSourceName)
	}
	driveri, ok := transport.GetDriver(u.Scheme)
	if !ok {
		return nil, errors.Errorf("cloudevent: unknown driver %q %s ", u.Scheme, u.String())
	}

	t, err = driveri.NewSinkTransport(ctx, u)
	if err != nil {
		return nil, errors.Errorf("cloudevent: failed to create sink %q %s ", u.Scheme, u.String())
	}

	return NewSink(ctx, t, opts...)
}
