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

package main

import (
	"context"
	"fmt"
	"github.com/tkeel-io/rule-util/stream"
	"time"
)

func ExampleKafka() {
	ctx := context.Background()
	testurl := "kafka://192.168.0.14:10002/topic-1/g2"
	sink, err := stream.OpenSink(testurl)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		msg := stream.NewMessage().
			SetPacketIdentifier(fmt.Sprintf("id-%d", i)).
			SetData([]byte("ABCD1234"))
		err := sink.Send(ctx, msg)
		fmt.Println(err)
	}

	time.Sleep(100 * time.Second)
	// Output:

}
