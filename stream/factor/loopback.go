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
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/types"
	"github.com/tkeel-io/rule-util/stream/utils"
	"net/url"
	"strings"
	"sync"
)

const DefaultQueueSize = 500

var (
	mutex  sync.Mutex
	caches = make(map[string]*loopback)
)

var _ types.Sink__ = (*loopback)(nil)
var _ types.Source__ = (*loopback)(nil)

type loopback struct {
	dns       string
	topic     string
	msgCh     chan Message
	running   bool
	callbacks []types.SourceHandle
	wg        *sync.WaitGroup
	lock      *sync.Mutex
}

func (this *loopback) String() string {
	return this.dns
}

func NewLoopback(ctx context.Context, dataSourceName string) (source *loopback) {
	return newLoopback(ctx, dataSourceName)
}

func NewLoopbackFromDNS(ctx context.Context, dataSourceName string) (source *loopback) {
	return newLoopback(ctx, dataSourceName)
}

func newLoopback(ctx context.Context, dataSourceName string) *loopback {
	var (
		u *url.URL
	)
	u, err := url.Parse(dataSourceName)
	if err != nil {
		return nil
	}
	if u.Scheme != "loopback" {
		return nil
	}
	if u.Host != "memqueue" {
		log.Error("Please use loopback://memqueue/{topic}/{group_id}[?param1=value1&paramN=valueN]")
		return nil
	}

	arr := strings.Split(u.Path, "/")
	if len(arr) != 3 && arr[0] != "" {
		log.Error("Please use loopback://memqueue/{topic}/{group_id}[?param1=value1&paramN=valueN]")
		return nil
	}
	topic := arr[1]

	mutex.Lock()
	defer mutex.Unlock()
	if topic == "" {
		return nil
	}
	if _, ok := caches[topic]; ok {
		return caches[topic]
	}
	caches[topic] = &loopback{
		dns:       dataSourceName,
		topic:     topic,
		msgCh:     make(chan Message, DefaultQueueSize),
		callbacks: make([]types.SourceHandle, 0),
		wg:        &sync.WaitGroup{},
		lock:      &sync.Mutex{},
	}
	caches[topic].wg.Add(1)
	return caches[topic]
}

func (this *loopback) Send(ctx context.Context, msg interface{}) error {
	if !this.running {
		if env := utils.Log.Check(log.InfoLevel, "loopback with none receiver"); env != nil {
			env.Write(logf.String("topic", this.topic))
		}
		return nil
	}
	for _, cb := range this.callbacks {
		if err := cb(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
func (this *loopback) SendAsync(ctx context.Context, event interface{}) (types.Promise, error) {
	promise := types.NewPromise()
	if !this.running {
		if env := utils.Log.Check(log.InfoLevel, "loopback with none receiver"); env != nil {
			env.Write(logf.String("topic", this.topic))
		}
		return promise.Finish(nil), nil
	}

	var err error = nil
	for _, cb := range this.callbacks {
		if err = cb(ctx, event); err != nil {
			break
		}
	}
	promise.Finish(err)
	return promise, nil
}

func (this *loopback) StartReceiver(ctx context.Context, fn types.SourceHandle) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.running = true
	this.callbacks = append(this.callbacks, fn)
	this.wg.Wait()
	return nil
}

func (this *loopback) Close(ctx context.Context) error {
	this.wg.Done()
	return nil
}
