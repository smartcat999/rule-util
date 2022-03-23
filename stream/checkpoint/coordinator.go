package checkpoint

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tkeel-io/rule-util/pkg/log"
	logf "github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/tkeel-io/rule-util/stream/types"
)

type pendingCheckpoint struct {
}

type CheckpointedSource interface {
	types.CheckpointedFunction
	types.CheckpointListener
}

type Coordinator struct {
	pending bool
	sources []CheckpointedSource
	sinks   []types.CheckpointedFunction
	ops     []types.CheckpointedFunction
	ticker  *time.Ticker //For trigger checkpoint only
	ctx     context.Context

	shutdown chan struct{}
	timeout  time.Duration
}

var coordinator *Coordinator
var lock sync.Mutex

func InitCoordinator(ctx context.Context, interval time.Duration) *Coordinator {
	lock.Lock()
	defer lock.Unlock()
	if coordinator == nil {
		coordinator = &Coordinator{
			pending:  false,
			sources:  make([]CheckpointedSource, 0),
			sinks:    make([]types.CheckpointedFunction, 0),
			ticker:   time.NewTicker(interval),
			ctx:      ctx,
			shutdown: make(chan struct{}, 0),
			timeout:  time.Second * 10,
		}
		coordinator.Run()
	} else {
		coordinator.ticker.Stop()
		coordinator.ticker = time.NewTicker(interval)
	}
	log.Debug(fmt.Sprintf("InitCoordinator"))
	return coordinator
}

func (c *Coordinator) Run() {
	go func() {
		for {
			select {
			case <-c.ticker.C:
				lock.Lock()
				if c.pending {
					lock.Unlock()
					return
				}
				c.pending = true
				lock.Unlock()
				now := time.Now().Unix()
				var checkpoint = NewCheckpointContext(now)
				wg := &sync.WaitGroup{}
				for _, r := range c.sources {
					wg.Add(1)
					go func(rr CheckpointedSource) {
						if err := rr.SnapshotState(checkpoint); err != nil {
							log.Info("Fail to trigger checkpoint",
								logf.Any("source", r),
								logf.Error(err))
							c.cancel(checkpoint)
						}
						wg.Done()
					}(r)
				}
				wg.Wait()

				for _, r := range c.ops {
					wg.Add(1)
					go func(rr types.CheckpointedFunction) {
						if err := rr.SnapshotState(checkpoint); err != nil {
							log.Info("Fail to trigger checkpoint",
								logf.Any("op", r),
								logf.Error(err))
							c.cancel(checkpoint)
						}
						wg.Done()
					}(r)
				}
				wg.Wait()

				for _, r := range c.sinks {
					wg.Add(1)
					go func(rr types.CheckpointedFunction) {
						if err := rr.SnapshotState(checkpoint); err != nil {
							log.Info("Fail to trigger checkpoint",
								logf.Any("sink", r),
								logf.Error(err))
							c.cancel(checkpoint)
						}
						wg.Done()
					}(r)
				}
				wg.Wait()

				for _, r := range c.sources {
					wg.Add(1)
					go func(rr CheckpointedSource) {
						if err := rr.NotifyCheckpointComplete(checkpoint.GetCheckpointId()); err != nil {
							log.Info("Fail to complete checkpoint",
								logf.Any("source", r),
								logf.Error(err))
							c.cancel(checkpoint)
						}
						wg.Done()
					}(r)
				}
				wg.Wait()

				lock.Lock()
				c.pending = false
				lock.Unlock()
			case <-c.ctx.Done():
				log.Info("Cancelling coordinator....")
				c.close()
				return
			case <-c.shutdown:
				log.Info("Cancelling coordinator....")
				c.close()
				return
			}
		}
	}()
}

func (c *Coordinator) cancel(snapshotContext types.FunctionSnapshotContext) {
	log.Error("cancel", logf.Any("snapshotContext", snapshotContext))
}

func (c *Coordinator) close() {
	if c.ticker != nil {
		c.ticker.Stop()
	}

}

func AddSource(source ...CheckpointedSource) {
	if coordinator == nil {
		log.Warn("coordinator nil", logf.Any("source", source))
		return
	}
	coordinator.sources = append(coordinator.sources, source...)
}

func AddSink(sink ...types.CheckpointedFunction) {
	if coordinator == nil {
		log.Warn("coordinator nil", logf.Any("sink", sink))
		return
	}
	coordinator.sinks = append(coordinator.sinks, sink...)
}

func AddOp(op ...types.CheckpointedFunction) {
	if coordinator == nil {
		log.Warn("coordinator nil", logf.Any("op", op))
		return
	}
	coordinator.ops = append(coordinator.sinks, op...)
}
