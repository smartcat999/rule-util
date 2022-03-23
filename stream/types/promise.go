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

package types

import "sync"

type promise struct {
	sync.Mutex
	err      error
	finished bool

	callbacks []func(err error)
}

func NewPromise() Promise {
	return new(promise)
}

func (p *promise) executeCallbacks() {
	if p.finished {
		return
	}
	for _, s := range p.callbacks {
		s(p.err)
	}
	p.finished = true
}

func (p *promise) Then(s func(err error)) Promise {
	p.Lock()
	defer p.Unlock()

	if p.finished {
		s(p.err)
	} else {
		p.callbacks = append(p.callbacks, s)
	}
	return p
}

func (p *promise) Finish(err error) Promise {
	p.Lock()
	defer p.Unlock()

	p.err = err

	p.executeCallbacks()
	return p
}
