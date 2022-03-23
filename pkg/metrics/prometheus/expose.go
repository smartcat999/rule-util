package prometheus

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	xregistry "github.com/tkeel-io/rule-util/pkg/registry"
	logf "github.com/tkeel-io/rule-util/pkg/logfield"
	"go.uber.org/zap"
)

func ExposedMetrics(reg *prometheus.Registry, conf *ExposedConf) error {

	var (
		err      error
		port     string
		listener net.Listener
	)

	if listener, err = net.Listen("tcp", conf.Addr+":0"); nil != err {
		logger.Bg().Error("listen metrics failed. error: %v", logf.Error(err))
		return err
	}

	//register etcd for metrics.
	if port = getPort(listener.Addr().String()); port == "" {
		return ErrResolvePort
	}

	xregistry.Register(context.Background(), conf.Etcd, conf.KV.Key(), conf.KV.Value(conf.Addr, port))

	servMux := http.NewServeMux()
	servMux.HandleFunc("/metrics",
		func(w http.ResponseWriter, r *http.Request) {
			promhttp.HandlerFor(reg,
				promhttp.HandlerOpts{}).ServeHTTP(w, r)
		})

	logger.Bg().Info("start prometheus http handler.",
		zap.String("ip", conf.Addr),
		zap.String("port", port))

	if err = http.Serve(listener, servMux); nil != err {
		logger.Bg().Fatal("start exporter failure.", logf.Error(err))
	}
	return err
}

func getPort(addr string) string {
	var port string
	arr := strings.Split(addr, ":")
	if len(arr) > 0 {
		port = arr[len(arr)-1]
		if "" != port {
			return port
		}
	}
	return ""
}
