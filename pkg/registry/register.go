package registry

import (
	"context"

	"github.com/tkeel-io/rule-util/pkg/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const heart_time = 10

var logger = log.GlobalLogger()

// type clease struct {
// 	leasse  clientv3.Lease
// 	leaseId clientv3.LeaseID
// }

// func newlease(ctx context.Context, endpoints []string, key string, value string) (*lease, error) {
// 	var (
// 		err     error
// 		client  *clientv3.Client
// 		lease   clientv3.Lease
// 		leaseID clientv3.LeaseID
// 	)
// 	if client, err = clientv3.New(clientv3.Config{Endpoints: endpoints}); err != nil {
// 		logger.Bg().Error("create etcd connection failed.",
// 			zap.Error(err))
// 		return nil, err
// 	}
// 	lease = clientv3.NewLease(client)
// 	if resp, err := lease.Grant(ctx, int64(heart_time)); err != nil {
// 		logger.Bg().Error("grant lease failed.",
// 			zap.Error(err))
// 		return nil, err
// 	} else {
// 		leaseID = resp.ID
// 		if _, err := client.Put(context.TODO(), key, value, clientv3.WithLease(leaseID)); err != nil {
// 			logger.Bg().Error("grpclb: set key with ttl to clientv3 failed.",
// 				zap.String("key", key),
// 				zap.Error(err))
// 			return nil, err
// 		} else {
// 			logger.Bg().Info("grpclb: register success.",
// 				zap.String("key", key),
// 				zap.String("value", value))
// 		}
// 	}
// 	return &lease{lease: lease, leaseId: leaseID}, nil
// }

// func (this *clease) renew() {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logger.Bg().Error("show: ", zap.Any("error", err))
// 		}
// 	}()

// 	var err error
// 	ctx, cancel := context.WithTimeout(context.Background(), (heart_time/3+1)*time.Second)
// 	defer cancel()
// 	if _, err = lease.KeepAliveOnce(ctx, leaseID); err != nil {
// 		logger.Bg().Error("grpclb: key KeepAlive failed.",
// 			zap.String("key", key),
// 			zap.Error(err))
// 	}
// 	return err
// }

func Register(ctx context.Context, endpoints []string, key string, value string) error {

	var (
		err     error
		client  *clientv3.Client
		lease   clientv3.Lease
		leaseID clientv3.LeaseID
	)
	if client, err = clientv3.New(clientv3.Config{Endpoints: endpoints}); err != nil {
		logger.Bg().Error("create etcd connection failed.",
			zap.Error(err))
		return err
	}
	lease = clientv3.NewLease(client)
	if resp, err := lease.Grant(ctx, int64(heart_time)); err != nil {
		logger.Bg().Error("grant lease failed.",
			zap.Error(err))
		return err
	} else {
		leaseID = resp.ID
		if _, err := client.Put(context.TODO(), key, value, clientv3.WithLease(leaseID)); err != nil {
			logger.Bg().Error("grpclb: set key with ttl to clientv3 failed.",
				zap.String("key", key),
				zap.Error(err))
			return err
		} else {
			logger.Bg().Info("grpclb: register success.",
				zap.String("key", key),
				zap.String("value", value))
		}
	}

	//keepalive
	if ch, err := lease.KeepAlive(context.Background(), leaseID); nil != err {
		return err
	} else {
		go func() {
			for {
				_ = <-ch
				//	log.Debug("keep alive", zap.Int64("ttl", ka.TTL))
			}
		}()
	}
	return nil

	// if err = renewLease(ctx, lease, leaseID, key); err != nil {
	// 	logger.Bg().Error("renew lease failed.", zap.Error(err))
	// 	return err
	// }

	// ticker := time.NewTicker((heart_time/3 + 1) * time.Second)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			if err = renewLease(ctx, lease, leaseID, key); err != nil {
	// 				logger.Bg().Error("renew lease failed.", zap.Error(err))
	// 			}
	// 		case <-ctx.Done():
	// 			ticker.Stop()
	// 			if _, err := client.Delete(context.Background(), key); err != nil {
	// 				logger.Bg().Error("grpclb: deregister failed.",
	// 					zap.String("key", key),
	// 					zap.Error(err))
	// 			}
	// 			return
	// 		}
	// 	}
	// }()
	// return nil
}

// func renewLease(ctx context.Context, lease clientv3.Lease, leaseID clientv3.LeaseID, key string) error {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logger.Bg().Error("show: ", zap.Any("error", err))
// 		}
// 	}()
// 	var err error
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	if _, err = lease.KeepAliveOnce(ctx, leaseID); err != nil {
// 		logger.Bg().Error("grpclb: key KeepAlive failed.",
// 			zap.String("key", key),
// 			zap.Error(err))
// 	}
// 	return err
// }
