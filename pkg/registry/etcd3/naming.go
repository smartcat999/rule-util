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
	"fmt"
	"time"

	"github.com/tkeel-io/rule-util/pkg/log"
	logf "github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/pkg/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

var (
	_ registry.Naming = &Discovery{}
)

type Discovery struct {
	client *clientv3.Client
	ticker *time.Ticker
	node   registry.Node
	logger log.Factory
}

func New(cfg *registry.Config) (*Discovery, error) {
	if cfg.Prefix != "" {
		log.GlobalLogger().Bg().Info("grpc: update prefix",
			logf.String("prefix", cfg.Prefix))
		global_prefix = cfg.Prefix
	}
	if len(cfg.Endpoints) == 0 {
		return nil, ErrNoAvailableEndpoints
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: 15 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	logger := log.GlobalLogger()
	d := &Discovery{
		client: client,
		logger: logger,
	}
	return d, nil
}

func (dc *Discovery) Register(ctx context.Context, srv *registry.Node, ttl int64) error {
	go func(srv *registry.Node, ttl int64) {
		ticker := time.NewTicker(time.Second * time.Duration(ttl))
	exit:
		for {
			path := dc.generatePath(srv)
			getResp, err := dc.client.Get(ctx, path)

			if err != nil {
				dc.logger.Bg().Error("register", zap.Error(err))
			}

			if getResp.Count == 0 {
				err = dc.withAlive(path, string(srv.Marshal()), ttl)
				if err != nil {
					dc.logger.Bg().Error("register", zap.Error(err))
				}
			}

			select {
			case <-ctx.Done():
				break exit
			case <-ticker.C:
			}
		}
	}(srv, ttl)
	return nil
}

func (dc *Discovery) withAlive(path string, metadata string, ttl int64) error {
	leaseResp, err := dc.client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	_, err = dc.client.Put(context.Background(), path, metadata, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	leaseId := int64(leaseResp.ID)
	ch, err := dc.client.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		return err
	} else {
		go func() {
			for {
				select {
				case resp := <-ch:
					if resp == nil {
						dc.logger.Bg().Info("lease: exit",
							logf.Int64("leaseId", leaseId),
							logf.Int64("ttl", ttl),
						)
						return
					} else {
						dc.logger.Bg().Info("lease: success",
							logf.Int64("leaseId", leaseId),
							logf.Int64("ttl", ttl),
						)
						goto END
					}
				}
				//			log.Debug("keep alive", zap.Int64("ttl", ka.TTL))
			END:
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}
	return nil
}

func (dc *Discovery) UnRegister(ctx context.Context, srv *registry.Node) error {
	_, err := dc.client.Delete(ctx, dc.generatePath(srv))
	return err
}

func (dc *Discovery) generatePath(srv *registry.Node) string {
	return fmt.Sprintf("/%s/%s/%s/%s",
		global_prefix, srv.AppID, registry_version, srv.Addr)

}
