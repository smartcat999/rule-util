package main

import (
	"context"
	"fmt"
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/stream/example"
	"github.com/tkeel-io/rule-util/stream/types"
	"time"
)

func receiveMessage(ctx context.Context, message interface{}) error {
	if msg, ok := message.(types.Message); ok {
		fmt.Println(string(msg.Data()))
	}
	return nil
}

func main() {
	log.InitLogger("test1", "test2")
	example.Co()
	qmqUrl := "kafka://127.0.0.1:9092/mdmp-topic/g1?logLevel=info"
	example.ReceiveMessage(qmqUrl, receiveMessage)
	qmqUrl = "kafka://127.0.0.1:9092/mdmp-topic/?logLevel=info"
	example.SendMessage(qmqUrl, 20)
	time.Sleep(20 * time.Second)
}
