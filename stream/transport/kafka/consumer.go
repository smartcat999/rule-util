package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/utils"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/Shopify/sarama"
)

type Consumer struct {
	lock     sync.Mutex
	shutdown chan struct{}

	sess        sarama.ConsumerGroup
	message     chan *sarama.ConsumerMessage
	topics      []string
	Addresses   []string
	checkpoint  map[types.CheckpointID]snapshot
	curSnapshot snapshot
	gs          sarama.ConsumerGroupSession
	option      *ConsumerOption
}

var _ types.SourceTransport = (*Consumer)(nil)

func (consumer *Consumer) NotifyCheckpointComplete(sid types.CheckpointID) error {
	offset := consumer.checkpoint[sid]
	for key, value := range offset {
		topic := consumer.topics[0]
		nextOffset := value.offset + 1
		consumer.gs.MarkOffset(topic, key, nextOffset, "test")
		if env := utils.Log.Check(log.DebugLevel, "NotifyCheckpointComplete:MarkOffset"); env != nil {
			env.Write(logf.String("topic", topic),
				logf.Int32("key", key),
				logf.Int64("offset", int64(value.offset)))
		}
	}
	consumer.lock.Lock()
	delete(consumer.checkpoint, sid)
	consumer.lock.Unlock()
	return nil
}

func (consumer *Consumer) SnapshotState(context types.FunctionSnapshotContext) error {
	sid := context.GetCheckpointId()
	consumer.checkpoint[sid] = make(snapshot)
	topic := consumer.topics[0]
	consumer.lock.Lock()
	curSnapshot := consumer.curSnapshot
	consumer.curSnapshot = make(snapshot)

	consumer.lock.Unlock()

	for key, value := range curSnapshot {
		consumer.checkpoint[sid][key] = value
		if env := utils.Log.Check(log.DebugLevel, "SnapshotState:Save"); env != nil {
			env.Write(logf.String("topic", topic),
				logf.Int32("key", key),
				logf.Int64("offset", int64(value.offset)))
		}
	}
	if env := utils.Log.Check(log.DebugLevel, "SnapshotState:Done"); env != nil {
		env.Write(logf.Any("sid", context.GetCheckpointId()))
	}
	return nil
}

func (consumer *Consumer) InitializeState(context types.FunctionInitializationContext) error {
	return nil
}

func (consumer *Consumer) Run(ctx context.Context, stx types.SourceContext) error {
	// Ensure group is not closed
	select {
	case <-consumer.shutdown:
		return ErrClosedConsumer
	default:
	}

	go func() {
		for {
			gerr, open := <-consumer.sess.Errors()
			if !open {
				return
			}
			if gerr != nil {
				//todo log & close
			}
		}
	}()

	go func() {
		sleepTime := time.Second
	groupLoop:
		for {
			ctx, cancel := context.WithCancel(context.Background())

			gerr := consumer.sess.Consume(ctx, consumer.topics, consumer)
			select {
			case <-ctx.Done():
				break groupLoop
			case <-consumer.shutdown:
				break groupLoop
			default:
			}

			cancel()
			if gerr != nil {
				if gerr != io.EOF {
					log.Error(fmt.Sprintf("KafkaCG group session error: %v\n", gerr))
				}
				if sleepTime < MaxReconnectTime {
					sleepTime = sleepTime * 2
				}
				time.Sleep(sleepTime)
				log.Info("Kafka consumer sleep")
			}
		}

		consumer.sess.Close()
	}()
	for {
		select {
		case tt := <-consumer.message:
			stx.Collect(tt.Value)
		case <-consumer.shutdown:
			break
		}
	}
	return nil
}

func (this *Consumer) parseOption(metadata map[string][]byte) (*ConsumerOption, error) {
	if s, ok := metadata["option"]; ok {
		opt := ConsumerOption{}
		err := json.Unmarshal([]byte(s), &opt)
		if err != nil {
			return nil, err
		}
		if opt.Addrs == "" {
			return nil, fmt.Errorf("empty addrs")
		}
		if opt.Topics == "" {
			return nil, fmt.Errorf("empty topics")
		}
		if opt.GroupID == "" {
			return nil, fmt.Errorf("empty group id")
		}
		return &opt, nil
	}
	return nil, fmt.Errorf("option not found")
}

func (consumer *Consumer) Open(ctx context.Context, configuration types.Configuration) error {
	metadata := configuration.Metadata()
	if opt, err := consumer.parseOption(metadata); err != nil {
		return err
	} else {
		consumer.option = opt
	}

	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	//config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 10 * time.Millisecond
	config.Consumer.Return.Errors = true
	config.Net.DialTimeout = time.Second

	config.ClientID = consumer.option.ClientID
	address := strings.Split(consumer.option.Addrs, ",")
	consumerGroup := consumer.option.GroupID
	consumer.topics = strings.Split(consumer.option.Topics, ",")

	group, err := sarama.NewConsumerGroup(address, consumerGroup, config)
	if err != nil {
		return err
	}
	consumer.sess = group
	consumer.message = make(chan *sarama.ConsumerMessage, 0)
	consumer.curSnapshot = make(snapshot)
	consumer.checkpoint = make(map[types.CheckpointID]snapshot)
	return nil
}

func (consumer *Consumer) Cancel(ctx context.Context) {
	panic("implement me")
}

func (consumer *Consumer) Close(ctx context.Context) error {
	consumer.lock.Lock()
	if consumer.shutdown == nil {
		consumer.shutdown = closedchan
	} else {
		close(consumer.shutdown)
	}
	consumer.lock.Unlock()
	if err := consumer.sess.Close(); err != nil {
		log.Error(err.Error())
	}
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	//	log.Println(session.MemberID())
	//	log.Println(session.GenerationID())
	if consumer.gs != session {
		consumer.gs = session
	}
	//	idx := rand.Int()
	//offsetMap := make(map[int32]int64)
	for message := range claim.Messages() {
		//		log.Printf("Message claimed[%d:%d:%d]: value = %s, timestamp = %v, topic = %s\n", idx, message.Partition, message.Offset, string(message.Value), message.Timestamp, message.Topic)
		consumer.message <- message
		consumer.lock.Lock()
		consumer.curSnapshot[message.Partition] = &ackContext{offset: message.Offset}
		consumer.lock.Unlock()
	}

	return nil
}
