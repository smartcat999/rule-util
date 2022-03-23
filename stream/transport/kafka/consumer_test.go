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

package kafka

import (
	"context"
	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/smartystreets/gunit"
	"log"
	"testing"
)

type Msg struct {
}

func (m Msg) Collect(msg interface{}) {
	log.Println(string(msg.([]byte)))
	//time.Sleep(1 * time.Second)
}

func (m Msg) MarkAsTemporarilyIdle() {
	panic("implement me")
}

func (m Msg) Context() context.Context {
	panic("implement me")
}

type Snap struct {
	id types.CheckpointID
}

func (s *Snap) IsRestored() bool {
	panic("implement me IsRestored")
}

func (s *Snap) GetStateStore() types.StateStore {
	panic("implement me GetStateStore")
}

func (s *Snap) GetCheckpointId() types.CheckpointID {
	return s.id
}

func (s *Snap) GetCheckpointTimestamp() types.CheckpointTimestamp {
	return types.CheckpointTimestamp(s.id)
}

func TestKafak(t *testing.T) {
	//gunit.RunSequential(new(KafakUnit), t)
	gunit.Run(new(KafakUnit), t)
}

type KafakUnit struct {
	*gunit.Fixture
}

func (this *KafakUnit) Setup() {
}

func (this *KafakUnit) TestBasic() {
}

func ExampleTest() {
	//URI := "kafka://192.168.0.14:10003/l1/g4"
	//
	//go func() {
	//	ctx := context.Background()
	//	sink, err := types.OpenSink(URI)
	//	idx := 0
	//	if err != nil {
	//		panic(err)
	//	}
	//	for {
	//		msg := stream.NewMessage().
	//			SetPacketIdentifier(fmt.Sprintf("id-%d", idx)).
	//			SetData([]byte("ABCD1234"))
	//		if err := sink.Send(ctx, msg); err != nil {
	//			log.Println(err)
	//		}
	//		time.Sleep(1 * time.Second)
	//		idx ++
	//	}
	//}()
	//
	//ctx := context.Background()
	//testurl, _ := url.Parse(URI)
	//consumer, err := (&Driver{}).NewSourceTransport(ctx, testurl)
	//log.Println(err)
	//msg := Msg{}
	//
	//log.Println("start")
	//go consumer.Run(ctx, msg)
	//snap := &Snap{
	//	types.CheckpointID(3),
	//}
	//log.Println("InitializeState")
	//consumer.InitializeState(nil)
	//
	//time.Sleep(3 * time.Second)
	//log.Println("SnapshotState", snap.id)
	//consumer.SnapshotState(snap)
	//time.Sleep(3 * time.Second)
	//log.Println("NotifyCheckpointComplete", snap.id)
	//consumer.NotifyCheckpointComplete(snap.id)
	//log.Println("notify", snap.id)
	//
	//time.Sleep(3 * time.Second)
	//log.Println("Close")
	//consumer.Close(ctx)
	//
	//time.Sleep(3 * time.Second)
	//
	//consumer2, err := (&Driver{}).NewSourceTransport(ctx, testurl)
	//log.Println("InitializeState", snap.id)
	//consumer2.InitializeState(snap)
	//log.Println("restart")
	//go consumer2.Run(ctx, msg)
	//time.Sleep(3 * time.Second)

	// Output:
	//
}
