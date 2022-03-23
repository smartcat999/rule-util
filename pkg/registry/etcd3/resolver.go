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
	"context"
	"encoding/json"
	"strings"

	"github.com/tkeel-io/rule-util/pkg/log"
	logf "github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/pkg/registry"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/resolver"
)

var (
	_ registry.Resolver = &Discovery{}
)

type NodeData struct {
	Addr     string
	Metadata map[string]string
}

type etcdResolver struct {
	rawAddr string
	cc      resolver.ClientConn
	cli     *clientv3.Client
	logger  log.Factory
}

func (dc *Discovery) GRPCResolver() resolver.Builder {
	return &etcdResolver{
		logger: dc.logger,
		cli:    dc.client,
	}
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	go r.watch("/" + global_prefix + "/" + target.Endpoint + "/" + registry_version + "/")
	return r, nil
}

func (r etcdResolver) Scheme() string {
	return schema
}

func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	//r.logger.Bg().Debug("ResolveNow") // TODO check
}

// Close closes the resolver.
func (r etcdResolver) Close() {
	r.logger.Bg().Info("Close")
}

func (r *etcdResolver) update(keyPrefix string) {

}
func (r *etcdResolver) watch(keyPrefix string) {
	var addrList []resolver.Address

	rsp, err := r.cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Info("etcd resolver watch",
			logf.Error(err))
	} else {
		for i := range rsp.Kvs {
			addr := string(rsp.Kvs[i].Key)
			addr = strings.Replace(addr, "[::]", "127.0.0.1", 1)
			nodeData := &NodeData{}
			err := json.Unmarshal([]byte(rsp.Kvs[i].Value), nodeData)
			if err == nil {
				addrList = append(addrList,
					resolver.Address{
						Addr:     strings.TrimPrefix(addr, keyPrefix),
						Metadata: &nodeData.Metadata,
					})
				r.logger.Bg().Debug(
					"etcd get",
					zap.String("key prefix", keyPrefix),
					zap.String("key", string(rsp.Kvs[i].Key)),
					zap.Any("addr list", addrList),
				)
			}
		}
	}

	//r.cc.NewAddress(addrList)
	r.cc.UpdateState(resolver.State{Addresses: addrList})

	rch := r.cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			r.logger.Bg().Debug("addr", zap.String("addr", addr))
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					nodeData := NodeData{}
					err := json.Unmarshal([]byte(ev.Kv.Value), &nodeData)
					if err == nil {
						addrList = append(addrList, resolver.Address{
							Addr:     addr,
							Metadata: &nodeData.Metadata,
						})
						//r.cc.NewAddress(addrList)
						r.cc.UpdateState(resolver.State{Addresses: addrList})
					} else {
						r.logger.Bg().Error(
							"json unmarshal",
							zap.Any("ev", ev),
							zap.Error(err),
						)
					}
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					//r.cc.NewAddress(addrList)
					r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
				r.logger.Bg().Debug(
					"etcd watch",
					zap.Any("ev", ev),
				)
			}
		}
	}
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
