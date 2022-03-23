package gops

import (
	"context"
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"github.com/google/gops/agent"
)

var shutdown = make(chan struct{})

// Init gops
func Run(ctx context.Context) {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal("gops agent listen",
			logf.Error(err))
	} else {
		log.Info("gops agent listen")
	}
	select {
	case <-shutdown:
	case <-ctx.Done():
	}
}
