package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/tkeel-io/rule-util/pkg/log"
	logf "github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/errors"
	"github.com/tkeel-io/rule-util/stream/types"
	"go.uber.org/atomic"
)

type FlushProducer interface {
	SendAsync(message *sarama.ProducerMessage, callback func(result *messagePacket)) (err error)
	Close() error
}

// mmTransport acts as both a http client and a http handler.
type sinkTransport struct {
	codec             types.Codec
	producer          FlushProducer
	pendingSleep      time.Duration
	pendingRecords    *atomic.Int32
	flushOnCheckpoint bool
	topic             string
	option            *SinkOption
	shutdown          chan struct{}
	logger            log.Factory
	tracer            opentracing.Tracer
	//eventCh     chan types.Message
}

// sinkTransport adheres to types.Transport.
var _ types.SinkAsyncTransport = (*sinkTransport)(nil)

func (s *sinkTransport) applyOptions(opts ...SinkOptionFn) error {
	for _, fn := range opts {
		if err := fn(s); err != nil {
			return err
		}
	}
	return nil
}

// New creates a new qmq producer transport.
func NewSinkTransport(opts ...SinkOptionFn) (*sinkTransport, error) {
	t := &sinkTransport{
		codec:             CodecV02{},
		shutdown:          make(chan struct{}, 0),
		pendingSleep:      100 * time.Millisecond,
		pendingRecords:    atomic.NewInt32(0),
		flushOnCheckpoint: false,
	}
	if err := t.applyOptions(opts...); err != nil {
		return nil, err
	}

	if t.logger == nil {
		logger := log.GlobalLogger()
		t.logger = logger
	}

	return t, checkSinkTransport(t)
}

func checkSinkTransport(t *sinkTransport) error {
	if t == nil {
		return errors.New("transport is nil")
	}
	if t.logger == nil {
		return errors.New("transport logger is nil")
	}
	return nil
}

func (s *sinkTransport) StartSpan(ctx context.Context, operationName string) opentracing.Span {
	// start span
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span = s.tracer.StartSpan(operationName, opentracing.ChildOf(span.Context()))
	} else {
		span = s.tracer.StartSpan(operationName)
	}
	return span
}

//InitializeState init consumer
//Do nothing, use remote queue offset
func (s *sinkTransport) InitializeState(context types.FunctionInitializationContext) error {
	return nil
}

func (s *sinkTransport) SnapshotState(context types.FunctionSnapshotContext) error {
	if s.flushOnCheckpoint {
		s.flush()
		if s.pendingRecords.Load() != 0 {
			return errors.Errorf("Pending record count must be zero at s point: %d", s.pendingRecords)
		}
	}
	return nil
}

func (s *sinkTransport) flush() {
	// The producer doesn't support flushing, we wait here
	// until all pending records are confirmed
	for s.pendingRecords.Load() > 0 {
		time.Sleep(s.pendingSleep)
	}
}

func (s *sinkTransport) AsyncInvoke(ctx context.Context, message interface{}, stx types.SinkContext) error {
	if stx == nil {
		stx = types.NewSinkContext()
	}
	promise := stx.Promise()
	topic := s.topic
	if topic == "" {
		err := errors.New("topic empty")
		s.logger.For(ctx).Error(
			"Send topic empty",
			logf.String("topic", topic),
			logf.Error(err),
		)
		return err
	}

	byt, err := s.Encode(message)
	if err != nil {
		return err
	}

	record := &sarama.ProducerMessage{
		Topic: topic,
		//Key:   sarama.StringEncoder(message),
		Value: sarama.ByteEncoder(byt),
	}

	if t, ok := message.(types.PublishMessage); ok {
		if t.Topic() != "" {
			record.Key = sarama.StringEncoder(t.Topic())
		} else {
			record.Key = sarama.StringEncoder(t.Entity())
		}
	} else if t, ok := message.(types.Message); ok {
		record.Key = sarama.StringEncoder(t.Entity())
	}

	if s.flushOnCheckpoint {
		s.pendingRecords.Add(1)
	}

	err = s.producer.SendAsync(record, func(result *messagePacket) {
		if s.flushOnCheckpoint {
			s.pendingRecords.Sub(1)
		}
		if result == nil || result.err != nil {
			rt := record.Topic
			rp := record.Partition
			s.logger.For(ctx).Error(
				"Send error",
				logf.String("topic", rt),
				logf.Int32("partition", rp),
				logf.Any("Data", record.Value),
			)
			err := fmt.Errorf("%s: %v [%v:%v:%v] ", errors.ErrSendError, result, rt, rp, record.Value)
			promise.Finish(err)
		} else {
			s.logger.For(ctx).Info("Send success",
				logf.String("Topic", topic))
			promise.Finish(nil)
		}
	})
	if err == nil {
		return nil
	}
	//if err == producer.ErrFullQueue {
	//	s.logger.Bg().Error("producer", logf.String("error", err.Error()))
	//	err = errors.ErrNeedSleep
	//} else {
	//	s.logger.Bg().Error("producer", logf.String("error", err.Error()))
	//	err = errors.ErrNeedRestart
	//}
	s.logger.Bg().Error("producer", logf.String("error", err.Error()))
	//err = errors.ErrNeedRestart
	return err
}

func (s *sinkTransport) Encode(message interface{}) ([]byte, error) {
	switch message := message.(type) {
	case []byte:
		body := message
		return body, nil
	default:
		body, err := s.codec.Encode(message)
		return body, err
	}
}

func (s *sinkTransport) Open(ctx context.Context, configuration types.Configuration) error {
	metadata := configuration.Metadata()
	if ret, err := s.parseOption(metadata); err != nil {
		return err
	} else {
		s.option = ret
	}
	s.option.config = NewSaramaConfig()
	s.option.retryCnt = SinkSendRetryCnt
	s.option.retryInterval = SinkSendRetryInterval
	s.option.config.Net.DialTimeout = 3 * time.Second

	topic := s.option.Topic

	c, err := NewProducer(s.option)
	if err != nil {
		return errors.Errorf("failed to create qmq producer")
	}
	s.producer = c

	s.topic = topic

	if err := checkSinkTransport(s); err != nil {
		return err
	}

	//s.InitializeState(nil)
	return nil

}

func (s *sinkTransport) Cancel(ctx context.Context) {
	s.Close(ctx)
	return
}

func (s *sinkTransport) Close(ctx context.Context) error {
	log.For(ctx).Info("transport to close")
	if s.producer != nil {
		s.producer.Close()
	}
	return nil
}

func (s *sinkTransport) parseOption(metadata map[string][]byte) (*SinkOption, error) {
	if s, ok := metadata["option"]; ok {
		opt := SinkOption{}
		err := json.Unmarshal([]byte(s), &opt)
		if err != nil {
			return nil, err
		}
		return &opt, nil
	}
	return nil, fmt.Errorf("option not found")
}
