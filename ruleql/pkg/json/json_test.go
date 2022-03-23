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

package json

import (
	"testing"

	"github.com/tkeel-io/rule-util/ruleql/pkg/json/gjson"
	"github.com/tkeel-io/rule-util/ruleql/pkg/json/jsonparser"
)

var json = `{
  "name": {"first": "Tom", "last": "Anderson", "age": 44},
  "age":37.1,
  "children": ["Sara","Alex","Jack"],
  "movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44},
    {"first": "Roger", "last": "Craig", "age": 68},
    {"first": "Jane", "last": "Murphy", "age": 47}
  ]
}`
var jsonByt = []byte(json)

func BenchmarkGet_gjson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.Get(json, "age")
	}
}

func BenchmarkGet_jsonparser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonparser.Get(jsonByt, "age")
	}
}

func BenchmarkGet_gjson2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.Get(json, "age")
		gjson.Get(json, "movie")
	}
}

func BenchmarkGet_jsonparser2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonparser.Get(jsonByt, "age")
		jsonparser.Get(jsonByt, "movie")
	}
}

func BenchmarkGet_gjson3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.Get(json, "age")
		gjson.Get(json, "age")
		gjson.Get(json, "age")
	}
}

func BenchmarkGet_jsonparser3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonparser.Get(jsonByt, "age")
		jsonparser.Get(jsonByt, "age")
		jsonparser.Get(jsonByt, "age")
	}
}

func BenchmarkGet_jsonparser_getunsafestring3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonparser.GetUnsafeString(jsonByt, "age")
		jsonparser.GetUnsafeString(jsonByt, "age")
		jsonparser.GetUnsafeString(jsonByt, "age")
	}
}
