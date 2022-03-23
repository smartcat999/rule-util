package prometheus

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type RateOpts prometheus.Opts

type RateGauge struct {
	data prometheus.Gauge
	opts RateOpts
}

func NewRateGauge(opts RateOpts) RateGauge {
	return RateGauge{
		opts: opts,
		data: prometheus.NewGauge(prometheus.GaugeOpts(opts)),
	}
}

func (this RateGauge) Describe(ch chan<- *prometheus.Desc) {
	ch <- this.desc()
}

func (this RateGauge) Collect(ch chan<- prometheus.Metric) {
	ch <- this.data
	this.reset()
}

func (this RateGauge) Add(num float64) {
	this.data.Add(num)
}

func (this RateGauge) Sub(num float64) {
	this.data.Sub(num)
}

func (this RateGauge) reset() {
	this.data.Set(0)
}

func (this RateGauge) desc() *prometheus.Desc {
	names := []string{}
	if this.opts.Namespace != "" {
		names = append(names, this.opts.Namespace)
	}
	if this.opts.Subsystem != "" {
		names = append(names, this.opts.Subsystem)
	}
	if this.opts.Name != "" {
		names = append(names, this.opts.Name)
	}
	return prometheus.NewDesc(strings.Join(names, "_"), this.opts.Help, nil, nil)
}
