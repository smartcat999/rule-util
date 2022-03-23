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

/*
 *  dataSourceName:
 *  qmq://[user[:password]@][net[(addr)]]/{topic}/{group_id}[?param1=value1&paramN=valueN]
 *  kafka://[user[:password]@][net[(addr)]]/{topic}/{group_id}/{partition_num}[?param1=value1&paramN=valueN]
 *
 */

package transport

import (
	"github.com/tkeel-io/rule-util/stream/types"
	"sort"
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]types.Driver)
)

func Register(name string, driver types.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("register transport is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("register called twice for driver " + name)
	}
	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	// For tests.
	drivers = make(map[string]types.Driver)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func GetDriver(Scheme string) (types.Driver, bool) {
	driversMu.RLock()
	driver, ok := drivers[Scheme]
	driversMu.RUnlock()
	return driver, ok
}
