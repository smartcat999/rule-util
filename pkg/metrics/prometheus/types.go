package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tkeel-io/rule-util/pkg/log"
	xregistry "github.com/tkeel-io/rule-util/pkg/registry"
)

type ExposedConf struct {
	Addr string
	Etcd []string
	KV   xregistry.Formater
}

type Option interface {
	Type() string
	Invoke(prometheus.Registerer) prometheus.Registerer
}

var logger = log.GlobalLogger()
