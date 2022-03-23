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
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/client"
	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/tkeel-io/rule-util/stream/utils"
)

const CloudEventsPrefix = "cloudevent"

// OpenSource
// cache 需要考虑每个slot都需要一个consume，先禁用
func OpenSource(dataSourceName string, opts ...client.Option) (types.SourceClient, error) {
	return OpenSourceContext(context.Background(), dataSourceName, opts...)
}

func OpenSourceContext(ctx context.Context, dataSourceName string, opts ...client.Option) (types.SourceClient, error) {
	//return cache.source(dataSourceName, opts...)
	if env := utils.Log.Check(log.DebugLevel, "source creat"); env != nil {
		env.Write(logf.String("dsn", dataSourceName))
	}
	source := NewLoopbackFromDNS(ctx, dataSourceName)
	if source != nil {
		return source, nil
	}
	return OpenStreamSourceContext(ctx, dataSourceName, opts...)
}

func OpenSink(dataSourceName string, opts ...client.Option) (types.SinkClient, error) {
	return OpenSinkContext(context.Background(), dataSourceName, opts...)
}

func OpenSinkContext(ctx context.Context, dataSourceName string, opts ...client.Option) (types.SinkClient, error) {
	//return cache.sink(dataSourceName, opts...)
	if env := utils.Log.Check(log.DebugLevel, "sink creat"); env != nil {
		env.Write(logf.String("dsn", dataSourceName))
	}
	sink := NewLoopbackFromDNS(ctx, dataSourceName)
	if sink != nil {
		return sink, nil
	}
	return OpenStreamSinkContext(ctx, dataSourceName, opts...)
}
