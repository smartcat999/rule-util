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

package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

var noopTrace = opentracing.NoopTracer{}

/*

Uses environment variables to set the tracer's Host
```
JAEGER_AGENT_HOST= 8.8.8.8
```
https://www.jaegertracing.io/docs/1.7/deployment/


1. Init logger and tracer

```
#cmd/serve.go:62
zapLogger := logger.With(zap.String("service", "youservicename"))
logger := log.NewFactory(zapLogger)
tracer := tracing.Init("youservicename", metricsFactory, logger)
opentracing.SetGlobalTracer(tracer)
```


3. Tag span
```
#internal/endpoint/event/endpoint.go:154
if span := opentracing.SpanFromContext(ctx); span != nil {
    // trace tag
    span.SetTag("topic", data.TopicName)
    // context with new span
    ctx = opentracing.ContextWithSpan(ctx, newSpan)
}
```


4. Trace log
```
#internal/endpoint/event/endpoint.go:159
s.logger.For(ctx).Info("Decode",
    zap.String("entity_id", entityID),
    zap.String("TopicName", data.TopicName),
    zap.ByteString("Payload", data.Payload))
```

*/

type TracingTag struct {
	key   string
	value interface{}
}

func NewSpan(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

func AddSpan(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		// start new span
		span := opentracing.StartSpan(operationName, opentracing.ChildOf(parentSpan.Context()))
		// context with new span
		ctx = opentracing.ContextWithSpan(ctx, span)
		return span, ctx
	}
	span := noopTrace.StartSpan(operationName)
	return span, ctx
}

func Tag(key string, value interface{}) TracingTag {
	return TracingTag{key, value}
}

func AddTag(ctx context.Context, tags ...TracingTag) context.Context {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		// start new span
		for _, tag := range tags {
			span.SetTag(tag.key, tag.value)
		}
		// context with new span
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	return ctx
}
