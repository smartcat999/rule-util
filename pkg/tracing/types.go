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

package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type Tracer = opentracing.Tracer

var GlobalTracer = opentracing.GlobalTracer
var SetGlobalTracer = opentracing.SetGlobalTracer

var String = log.String
var Int = log.Int
var Int32 = log.Int32
var Int64 = log.Int64
var Float32 = log.Float32
var Float64 = log.Float64
var Bool = log.Bool
var Error = log.Error
var Any = log.Object
var Object = log.Object
