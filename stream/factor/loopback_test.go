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

package factor

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

func ExampleTestLoopback() {
	sink, _ := OpenSink("loopback://memqueue/abc1/group1")
	source, _ := OpenSource("loopback://memqueue/abc1/group1")
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	outbox := []Message{
		NewMessage().SetPacketIdentifier("message-1"),
		NewMessage().SetPacketIdentifier("message-2"),
		NewMessage().SetPacketIdentifier("message-3"),
		NewMessage().SetPacketIdentifier("message-4"),
	}
	inbox := []Message{}
	go source.StartReceiver(ctx, func(ctx context.Context, message Message) error {
		inbox = append(inbox, message)
		fmt.Println("Receiver...", message.PacketIdentifier())
		wg.Done()
		return nil
	})
	time.Sleep(1 * time.Millisecond)
	for idx, message := range outbox {
		wg.Add(1)
		if idx%2 == 0 {
			fmt.Println("Send", message.PacketIdentifier())
			sink.SendAsync(ctx, message)
		} else {
			fmt.Println("SendAsync", message.PacketIdentifier())
			if promise, err := sink.SendAsync(ctx, message); err != nil {
				fmt.Println("SendAsync Error...", message.PacketIdentifier(), err)
			} else {
				promise.Then(func(err error) {
					fmt.Println("Promise...", message.PacketIdentifier())
				})
			}
		}
		fmt.Println("Send Done...", message.PacketIdentifier())
	}

	fmt.Println("Wait...")
	wg.Wait()

	for i := 0; i < len(outbox); i++ {
		if outbox[i] != inbox[i] {
			fmt.Println("err", outbox[i], inbox[i])
		}
	}
	fmt.Println("Done")

	//Output:
	//Send message-1
	//Receiver... message-1
	//Send Done... message-1
	//SendAsync message-2
	//Receiver... message-2
	//Promise... message-2
	//Send Done... message-2
	//Send message-3
	//Receiver... message-3
	//Send Done... message-3
	//SendAsync message-4
	//Receiver... message-4
	//Promise... message-4
	//Send Done... message-4
	//Wait...
	//Done
}

func ExampleTestTwoLoopback() {
	wg := &sync.WaitGroup{}
	ctx := context.Background()
	inbox1 := []Message{}
	inbox2 := []Message{}
	sink1, _ := OpenSink("loopback://memqueue/abc1/group1")
	source1, _ := OpenSource("loopback://memqueue/abc1/group1")
	sink2, _ := OpenSink("loopback://memqueue/abc2/group1")
	source2, _ := OpenSource("loopback://memqueue/abc2/group1")
	outbox := []Message{
		NewMessage().SetPacketIdentifier("message-1"),
		NewMessage().SetPacketIdentifier("message-2"),
		NewMessage().SetPacketIdentifier("message-3"),
		NewMessage().SetPacketIdentifier("message-4"),
	}
	wg.Add(len(outbox))
	go source1.StartReceiver(ctx, func(ctx context.Context, message Message) error {
		inbox1 = append(inbox1, message)
		log.Print("inbox1 receiver...", message.PacketIdentifier())
		wg.Done()
		return nil
	})
	go source2.StartReceiver(ctx, func(ctx context.Context, message Message) error {
		inbox2 = append(inbox2, message)
		log.Print("inbox2 receiver...", message.PacketIdentifier())
		return nil
	})

	time.Sleep(time.Second * 1)
	go func() {
		for _, message := range outbox {
			log.Print("SendAsync...", message.PacketIdentifier())
			_, _ = sink1.SendAsync(ctx, message)
			_, _ = sink2.SendAsync(ctx, message)
			//fmt.Println(p,e)
		}
	}()

	fmt.Println("Wait...")
	//wg.Wait()
	time.Sleep(time.Second * 1)

	fmt.Println("inbox1...")
	for i := 0; i < len(inbox1); i++ {
		fmt.Println(inbox1[i].PacketIdentifier())
	}
	fmt.Println("inbox2...")
	for i := 0; i < len(inbox2); i++ {
		fmt.Println(inbox2[i].PacketIdentifier())
	}
	fmt.Println("Done")

	//Output:
	//
	//Wait...
	//inbox1...
	//message-1
	//message-2
	//message-3
	//message-4
	//inbox2...
	//message-1
	//message-2
	//message-3
	//message-4
	//Done
}
