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

package pprof

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
)

var profiles map[string]string

func StartCPUProfile(name string) {
	filename := fmt.Sprintf("cpu_%s.profile", name)
	cpuf, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuf)
}

func StopCPUProfile() {
	pprof.StopCPUProfile()
}

func StartPrefTrace(name string) {
	filename := fmt.Sprintf("trace_%s.profile", name)
	cput, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	trace.Start(cput)
}

func StopPrefTrace() {
	trace.Stop()
}

func WriteHeapProfile(name string) {
	filename := fmt.Sprintf("mem_%s.profile", name)
	memf, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	if err := pprof.WriteHeapProfile(memf); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	memf.Close()
}
