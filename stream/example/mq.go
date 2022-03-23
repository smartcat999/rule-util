package example

import (
	"context"
	"fmt"
	"github.com/tkeel-io/rule-util/stream"
	"github.com/tkeel-io/rule-util/stream/checkpoint"
	"github.com/tkeel-io/rule-util/stream/types"
	"time"
)

func SendMessage(mqUrl string, msgCnt int) {
	sink, err := stream.OpenSink(mqUrl)
	if err != nil {
		fmt.Println(err)
	}
	ctx := context.Background()
	for i := 0; i < msgCnt; i++ {
		msg := stream.NewMessage()
		msg.SetData([]byte(fmt.Sprintf("%d", i)))
		err = sink.Send(ctx, msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ReceiveMessage(mqUrl string, receiveMessage types.SourceHandle) {
	source, err := stream.OpenSource(mqUrl)
	if err != nil {
		fmt.Println(err)
	}
	if err = source.StartReceiver(context.Background(), receiveMessage); err != nil {
		fmt.Println(err)
	}
}

func Co() {
	co := checkpoint.InitCoordinator(context.Background(), 10*time.Second)
	co.Run()
}
