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

package etcdv3

import (
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"
)

const schema = "etcd3"
const registry_version = "v1.0.0"

//var global_prefix = "etcd3"
var global_prefix = "grpc-lb"

var (
	ErrNoAvailableEndpoints = errors.New("grpc: no available endpoints")
	ErrNoAvailablePrefix    = errors.New("grpc: no available prefix")
)
var tracer opentracing.Tracer
var metricsFactory = metrics.NullFactory
