package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

func Setup(reg prometheus.Registerer, opts ...Option) prometheus.Registerer {

	for _, opt := range opts {
		reg = opt.Invoke(reg)
	}
	return reg
}
