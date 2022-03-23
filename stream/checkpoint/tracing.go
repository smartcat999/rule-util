package checkpoint

import (
	"context"

	v1 "github.com/tkeel-io/rule-util/metadata/v1"

	//"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

var (
	mdmpTag = opentracing.Tag{Key: string(ext.Component), Value: "MDMP-tag"}
)

func SpanFromEvent(event interface{}, methodName string) (span opentracing.Span) {
	var (
		err           error
		parentSpanCtx opentracing.SpanContext
	)
	tracer := opentracing.GlobalTracer()
	message, ok := event.(v1.Attribution)
	if !ok {
		//fmt.Println("err == opentracing.ErrSpanContextNotFound")
		opts := []opentracing.StartSpanOption{
			mdmpTag,
		}
		span = tracer.StartSpan(methodName, opts...)
		return span
	}

	parentSpanCtx, err = tracer.Extract(opentracing.TextMap, CloudeventsCarrier{message})
	if err == opentracing.ErrSpanContextNotFound {
		//fmt.Println("err == opentracing.ErrSpanContextNotFound")
		opts := []opentracing.StartSpanOption{
			mdmpTag,
		}
		span = tracer.StartSpan(methodName, opts...)
		return span
	}

	if err != nil {
		return nil
	} else {
		//fmt.Println("opentracing.ChildOf(parentSpanCtx)")
		opts := []opentracing.StartSpanOption{
			opentracing.ChildOf(parentSpanCtx),
			mdmpTag,
		}
		span = tracer.StartSpan(methodName, opts...)
	}
	return span

}

func EventWithSpan(event interface{}, span opentracing.Span) interface{} {
	if message, ok := event.(v1.Attribution); ok {
		tracer := opentracing.GlobalTracer()
		_ = tracer.Inject(span.Context(), opentracing.TextMap, CloudeventsCarrier{message})
	}
	return event
}

//StartSpan
func StartSpan(methodName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return opentracing.StartSpan(methodName, opts...)
}

//StartSpanFromContext
func StartSpanFromContext(ctx context.Context, methodName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	tracer := opentracing.GlobalTracer()
	return opentracing.StartSpanFromContextWithTracer(ctx, tracer, methodName, opts...)
}

//CloudeventsCarrier
type CloudeventsCarrier struct {
	Attribution
}

// implement EventReader
func (c CloudeventsCarrier) ForeachKey(handler func(key, val string) error) error {
	return c.Attribution.ForeachAttr(handler)
}

// implement EventWriter
func (c CloudeventsCarrier) Set(key, val string) {
	//fmt.Println("[+]Set", key, val)
	c.Attribution.SetAttr(key, []byte(val))
}
