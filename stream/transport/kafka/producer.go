package kafka

import (
	"errors"
	"fmt"
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/tkeel-io/rule-util/stream/utils"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type Callback func()

type messagePacket struct {
	message  *sarama.ProducerMessage
	callback func(result *messagePacket)
	err      error
}
type ProducerMessage struct {
	Topic string
	Key   sarama.StringEncoder
	Value []byte
}

// Producer abstracts the kafka producer
type Producer interface {
	// Emit sends a message to topic.
	Emit(message *sarama.ProducerMessage) types.Promise
	Close() error
}

type producer struct {
	producer sarama.AsyncProducer
	stop     chan bool
	done     chan bool
	wg       sync.WaitGroup
}

// NewProducer creates new kafka producer for passed brokers.
func newProducer(option *SinkOption) (Producer, error) {
	brokers := option.Addrs
	config := option.config
	aprod, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to start Sarama producer: %v", err)
	}

	p := producer{
		producer: aprod,
		stop:     make(chan bool),
		done:     make(chan bool),
	}

	go p.run()

	return &p, nil
}

// Close stops the producer and waits for the Success/Error channels to drain.
// Emitting to a closing/closed producer results in write-to-closed-channel panic
func (p *producer) Close() error {
	// do an async close to get the rest of the success/error messages to avoid
	// leaving unfinished promises.
	p.producer.AsyncClose()

	// wait for the channels to drain
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.NewTimer(60 * time.Second).C:
	}

	return nil
}

// Emit emits a key-value pair to topic and returns a Promise that
// can be checked for errors asynchronously
func (p *producer) Emit(message *sarama.ProducerMessage) types.Promise {
	promise := types.NewPromise()
	message.Metadata = promise
	p.producer.Input() <- message
	//p.producer.Input() <- &sarama.ProducerMessage{
	//	Topic:    message.Topic,
	//	Key:      sarama.StringEncoder(message.Key),
	//	Value:    sarama.ByteEncoder(message.Value),
	//	Metadata: promise,
	//}
	return promise
}

// resolve or reject a promise in the message's metadata on Success or Error
func (p *producer) run() {
	p.wg.Add(2)
	go func() {
		defer p.wg.Done()
		for {
			err, ok := <-p.producer.Errors()

			// channel closed, the producer is stopping
			if !ok {
				return
			}
			err.Msg.Metadata.(types.Promise).Finish(err.Err)
		}
	}()

	go func() {
		defer p.wg.Done()
		for {
			msg, ok := <-p.producer.Successes()
			// channel closed, the producer is stopping
			if !ok {
				return
			}
			msg.Metadata.(types.Promise).Finish(nil)
		}
	}()
	p.wg.Wait()
	fmt.Println("Exit~~~~")
}

type flushProducer struct {
	producer  Producer
	wg        *sync.WaitGroup
	lock      sync.Mutex
	messageCh chan *messagePacket
	retryCh   chan *messagePacket
	maxRetry  int
	interval  time.Duration
}

func (p *flushProducer) Close() error {
	return p.producer.Close()
}

func (p *flushProducer) SnapshotState() error {
	//wg := p.wg
	//wg.Done()
	panic("TODO....")
}

func (p *flushProducer) SendAsync(message *sarama.ProducerMessage, callback func(result *messagePacket)) error {
	p.producer.Emit(message).Then(
		func(err error) {
			if err != nil {
				p.doRetry(&messagePacket{message, callback, nil})
			} else {
				callback(&messagePacket{message, callback, nil})
			}
		})
	return nil
}

func (p *flushProducer) doRetry(mm *messagePacket) {
	go func() {
		var wg = &sync.WaitGroup{}
		var result error
		var retryCnt = p.maxRetry
		for {
			if env := utils.Log.Check(log.DebugLevel, "doRetry"); env != nil {
				env.Write(logf.Any("message", mm.message),
					logf.Int("retryStatus", retryCnt))
			}
			if retryCnt == 0 {
				err := errors.New("too many retry")
				mm.err = err
				mm.callback(mm)
				return
			}
			wg.Add(1)
			promise := p.producer.Emit(mm.message)
			promise.Then(func(err error) {
				result = err
				wg.Done()
			})
			wg.Wait()
			if result != nil {
				retryCnt--
			} else {
				mm.callback(mm)
				return
			}
			time.Sleep(p.interval)
		}
	}()
}

func NewProducer(option *SinkOption) (fpd *flushProducer, err error) {
	var pd Producer

	if pd, err = newProducer(option); err == nil {
		fpd = &flushProducer{
			producer: pd,
			maxRetry: option.retryCnt,
			interval: option.retryInterval,
		}
	}

	return
}
