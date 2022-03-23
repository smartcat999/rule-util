package kafka

import (
	"context"
	"fmt"
	"github.com/tkeel-io/rule-util/stream/errors"
	"github.com/tkeel-io/rule-util/stream/transport"
	"github.com/tkeel-io/rule-util/stream/types"
	"net/url"
	"strings"
)

type Driver struct {
}

func (*Driver) NewSourceTransport(ctx context.Context, uri *url.URL) (types.SourceTransport, error) {
	arr := strings.Split(uri.Path, "/")
	if len(arr) != 3 || arr[0] != "" {
		fmt.Println("kafka://[user[:password]@][net[(addr)]]/{topic}/{group_id}[?partition=partition&offset=offset]")
		return nil, fmt.Errorf("qmq: parse error ")
	}
	consumer := Consumer{
		shutdown: make(chan struct{}, 0),
	}
	if cfg, err := NewConsumerConfig(map[string]interface{}{"option": &ConsumerOption{
		Addrs:    uri.Host,
		GroupID:  arr[2],
		Topics:   arr[1],
		ClientID: "sarama",
	},}); err != nil {
		return nil, err
	} else if err := consumer.Open(ctx, cfg); err != nil {
		return nil, err
	}
	return &consumer, nil
}

// NewSink  qmq://[user[:password]@][net[(addr)]]/{topic}/{group_id}[?param1=value1&paramN=valueN]
// @TODO support more options
func (m *Driver) NewSinkTransport(ctx context.Context, u *url.URL) (types.SinkAsyncTransport, error) {
	arr := strings.Split(u.Path, "/")
	if len(arr) != 3 && arr[0] != "" {
		fmt.Println("kafka://[user[:password]@][net[(addr)]]/{topic}/{group_id}[?param1=value1&paramN=valueN]")
		return nil, errors.Errorf("kafka: parse error ")
	}
	c, err := NewSinkTransport()
	if err != nil {
		return nil, errors.Errorf("failed to create kafka consumer")
	}

	logLevel := u.Query().Get("logLevel")
	if logLevel == "" {
		logLevel = "error"
	}
	if cfg, err := NewProducerConfig(map[string]interface{}{"option": &SinkOption{
		Addrs:    strings.Split(u.Host, ","),
		LogLevel: logLevel,
		Topic:    arr[1],
		GroupID:  arr[2],
	},}); err != nil {
		return nil, err
	} else if err := c.Open(ctx, cfg); err != nil {
		return nil, err
	}
	return c, nil
}

func init() {
	transport.Register("kafka", &Driver{})
}
