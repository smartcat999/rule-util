package client

import (
	"context"
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/errors"
	"github.com/tkeel-io/rule-util/stream/types"
	"sync"
	"time"
)

// New produces a new client with the provided source object and applied
// client options.
func NewSource(ctx context.Context, t types.SourceTransport, opts ...Option) (types.SourceClient, error) {
	c := &ceClient{
		ctx:            ctx,
		source:         t,
		codec:          CodecV02{},
		commitInterval: 1000 * time.Millisecond,
	}
	if err := c.applyOptions(opts...); err != nil {
		return nil, err
	}
	return c, nil
}

// New produces a new client with the provided source object and applied
// client options.
func NewSink(ctx context.Context, t types.SinkAsyncTransport, opts ...Option) (types.SinkClient, error) {
	c := &ceClient{
		ctx:   ctx,
		sink:  t,
		codec: CodecV02{},
	}
	if err := c.applyOptions(opts...); err != nil {
		return nil, err
	}
	return c, nil
}

type ceClient struct {
	source types.SourceTransport
	sink   types.SinkAsyncTransport
	fn     types.SourceHandle

	receiverMu     sync.Mutex
	once           sync.Once
	codec          Codec
	ctx            context.Context
	shutdown       chan struct{}
	lock           sync.Mutex
	commitInterval time.Duration

	//eventDefaulterFns []EventDefaulter
}

func (c *ceClient) SendAsync(ctx context.Context, event interface{}) (types.Promise, error) {
	return c.obsSendAsync(ctx, event)
}

func (c *ceClient) Close(ctx context.Context) error {
	c.lock.Lock()
	if c.shutdown == nil {
		c.shutdown = closedchan
	} else {
		close(c.shutdown)
	}
	c.lock.Unlock()
	if c.sink != nil && c.source != nil {
		panic("source and a sink repel each other")
	}
	if c.sink != nil {
		return c.sink.Close(ctx)
	}
	if c.source != nil {
		return c.source.Close(ctx)
	}
	return nil
}

func (c *ceClient) obsSendAsync(ctx context.Context, event interface{}) (types.Promise, error) {
	// Confirm we have a source set.
	if c.sink == nil {
		err := errors.Errorf("client not ready, sink not initialized")
		return types.NewPromise().Finish(nil), err
	}

	// Send the event over the source.
	sc := types.NewSinkContext()

	return sc.Promise(), c.sink.AsyncInvoke(ctx, event, sc)
}

// Receive is called from from the source on event delivery.
func (c *ceClient) Receive(ctx context.Context, event types.Message) error {
	err := c.obsReceive(ctx, event)
	return err
}

func (c *ceClient) obsReceive(ctx context.Context, event types.Message) error {
	if c.fn != nil {
		err := c.fn(ctx, event)
		return err
	}
	return nil
}

func (c *ceClient) applyOptions(opts ...Option) error {
	for _, fn := range opts {
		if err := fn(c); err != nil {
			return err
		}
	}
	return nil
}

// StartReceiver sets up the given fn to handle Receive.
// See Client.StartReceiver for details. This is a blocking call.
func (c *ceClient) StartReceiver(ctx context.Context, fn types.SourceHandle) error {
	c.receiverMu.Lock()
	defer c.receiverMu.Unlock()

	if c.source == nil {
		return errors.Errorf("client not ready, source not initialized")
	}

	c.fn = fn
	defer func() {
		c.fn = nil
	}()

	go func() {
		for {
			select {
			case <-c.shutdown:
				return
			default:
			}
			time.Sleep(c.commitInterval)
			c.commit()
		}
	}()

	return c.source.Run(ctx, c)

}

func (c *ceClient) Collect(msg interface{}) {
	var (
		byt   []byte
		ok    bool
		err   error
		event types.Message
	)
	if byt, ok = msg.([]byte); !ok {
		return
	}
	if event, err = c.codec.Decode(byt); err != nil {
		log.For(c.ctx).Error("codec.Decode",
			logf.Error(errors.ErrEncodeError))
		return
	}
	if event.Version() != types.MESSAGE_VERSION_1_0 {
		log.For(c.ctx).Error("codec.Decode",
			logf.Error(errors.ErrEncodeError))
		return
	}
	c.fn(c.ctx, event)
}

func (c *ceClient) MarkAsTemporarilyIdle() {

}

func (c *ceClient) commit() {
	ss := NewSnapContext()
	c.source.SnapshotState(ss)
	c.source.NotifyCheckpointComplete(ss.GetCheckpointId())
}
