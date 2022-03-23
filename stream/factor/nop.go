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

package factor

import (
	"context"
	"fmt"
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/tkeel-io/rule-util/stream/utils"
	"net/url"
)

var NoopSink = NewNoopStream(context.Background(), "nop://null/")

var _ types.Sink__ = (*loopback)(nil)
var _ types.Source__ = (*loopback)(nil)

type noopStream struct {
	logger         log.Factory
	logged         bool
	dataSourceName string
}

var _ types.SinkClient = (*noopStream)(nil)
var _ types.SourceClient = (*noopStream)(nil)

func NewNoopStream(ctx context.Context, dataSourceName string) (source *noopStream) {
	return newNoopStream(ctx, dataSourceName)
}

func NewNoopStreamFromDNS(ctx context.Context, dataSourceName string) (source *noopStream) {
	return newNoopStream(ctx, dataSourceName)
}

func newNoopStream(ctx context.Context, dataSourceName string) *noopStream {
	var (
		u      *url.URL
		logged bool
	)
	u, err := url.Parse(dataSourceName)
	if err != nil {
	}
	if u.Scheme != "nop" {
		return nil
	}
	switch u.Host {
	case "null":
		logged = false
	case "stdio":
		logged = true
	default:
		logged = false
		log.Error("Please use nop://null/{topic}/{group_id}[?param1=value1&paramN=valueN]")
		return nil
	}

	return &noopStream{
		dataSourceName: dataSourceName,
		logged:         logged,
		logger:         log.GlobalLogger(),
	}
}

func (this *noopStream) String() string {
	return "noop://null/"
}

func (this *noopStream) Send(ctx context.Context, msg interface{}) (error) {
	if env := utils.Log.Check(log.DebugLevel, "noop send"); env != nil {
		env.Write(logf.Any("dataSourceName", this.dataSourceName))
	}

	if this.logged {
		fmt.Println(msg)
	}
	return nil
}

func (this *noopStream) StartReceiver(ctx context.Context, fn types.SourceHandle) error {
	this.logger.For(ctx).Warn("listen on noop stream")
	return nil
}

func (this *noopStream) SendAsync(ctx context.Context, event interface{}) (types.Promise, error) {
	if env := utils.Log.Check(log.DebugLevel, "noop send"); env != nil {
		env.Write(logf.Any("dataSourceName", this.dataSourceName))
	}

	if this.logged {
		fmt.Println(event)
	}
	return types.NewPromise().Finish(nil), nil
}

func (this *noopStream) Close(ctx context.Context) error {
	return nil
}
