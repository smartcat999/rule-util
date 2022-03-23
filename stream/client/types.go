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

package client

import (
	"github.com/tkeel-io/rule-util/stream/types"
	"time"
)

type Codec interface {
	Encode(types.Message) ([]byte, error)
	Decode([]byte) (types.Message, error)
}

type snapContext struct {
	id types.CheckpointID
}

func NewSnapContext() *snapContext {
	return &snapContext{
		types.CheckpointID(time.Now().Unix()),
	}
}

func (s *snapContext) IsRestored() bool { return false }

func (s *snapContext) GetStateStore() types.StateStore { return nil }

func (s *snapContext) GetCheckpointId() types.CheckpointID {
	return s.id
}

func (s *snapContext) GetCheckpointTimestamp() types.CheckpointTimestamp {
	return types.CheckpointTimestamp(s.id)
}

// closedchan is a reusable closed channel.
var closedchan = make(chan struct{})

func init() {
	close(closedchan)
}
