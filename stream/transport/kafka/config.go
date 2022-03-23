package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/Shopify/sarama"
	"time"
)

const (
	// producer flush configuration
	defaultFlushFrequency     = 100 * time.Millisecond
	defaultFlushBytes         = 64 * 1024
	defaultProducerMaxRetries = 10
	defaultMaxProcessingTime  = 3 * time.Second
)

func NewConfig(cfg map[string]interface{}) (types.Configuration, error) {
	config := &configuration{
		name:     "",
		metadata: make(map[string][]byte),
	}
	for key, obj := range cfg {
		if byt, err := json.Marshal(obj); err != nil {
			return nil, fmt.Errorf("config marshal fail")
		} else {
			config.metadata[key] = byt
		}
	}
	return config, nil
}

// NewConfig creates a (bsm) sarama configuration with default values.
func NewSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_1_0

	// consumer configuration
	config.Consumer.Return.Errors = true
	config.Consumer.MaxProcessingTime = defaultMaxProcessingTime
	// this configures the initial offset for streams. Tables are always
	// consumed from OffsetOldest.
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// producer configuration
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = defaultFlushFrequency
	config.Producer.Flush.Bytes = defaultFlushBytes
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = defaultProducerMaxRetries

	return config
}

type configuration struct {
	name     string
	metadata map[string][]byte
}

func NewProducerConfig(cfg map[string]interface{}) (types.Configuration, error) {
	config := &configuration{
		name:     "",
		metadata: make(map[string][]byte),
	}
	for key, obj := range cfg {
		if byt, err := json.Marshal(obj); err != nil {
			return nil, fmt.Errorf("config marshal fail")
		} else {
			config.metadata[key] = byt
		}
	}
	return config, nil
}

func NewConsumerConfig(cfg map[string]interface{}) (types.Configuration, error) {
	config := &configuration{
		name:     "",
		metadata: make(map[string][]byte),
	}
	for key, obj := range cfg {
		if byt, err := json.Marshal(obj); err != nil {
			return nil, fmt.Errorf("config marshal fail")
		} else {
			config.metadata[key] = byt
		}
	}
	return config, nil
}

func (c *configuration) Name() string {
	return "config"
}

func (c *configuration) Metadata() map[string][]byte {
	return c.metadata
}

