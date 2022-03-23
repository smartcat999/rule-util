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

package registry

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/resolver"
)

// Config discovery configures.
type Config struct {
	Prefix    string
	Endpoints []string
}

type Node struct {
	// Hostname hostname.
	Hostname string `json:"hostname"`
	// AppID service id.
	AppID string `json:"appid"`
	// Addrs app instance address
	// format: scheme://host
	Addr            string `json:"addr"`
	Version         string `json:"version"`
	LastTimestampMS int64  `json:"last_timestamp"`
	// Metadata extension field.
	Metadata map[string]string `json:"metadata"`
}

type Naming interface {
	Register(context.Context, *Node, int64) error
	UnRegister(context.Context, *Node) error
}

type Resolver interface {
	GRPCResolver() resolver.Builder
}

func (n *Node) Marshal() []byte {
	bt, err := json.Marshal(n)
	if err == nil {
		return bt
	}
	return []byte{}
}

type Formater interface {
	Key() string
	Value(args ...string) string
}
