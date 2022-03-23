package debug

import (
	"fmt"
	"github.com/tkeel-io/rule-util/pkg/log"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
)

type HandleFunc struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

// Init pprof
// endpoints 127.0.0.1:7001,10.0.0.1:7001
// endpoints 0.0.0.0:7001
func Init(endpoints []string, handles ...*HandleFunc) {
	logger := log.GlobalLogger()
	pprofServeMux := http.NewServeMux()

	for _, handle := range handles {
		prefix := "/debug/" + handle.Pattern
		prefix = strings.ReplaceAll(prefix, "//", "/")
		pprofServeMux.HandleFunc(prefix, handle.Handler)
		logger.Bg().Info(
			"handel",
			zap.String("prefix", prefix),
		)
	}

	if len(endpoints) == 0 {
		endpoints = append(endpoints, ":0")
	}

	for _, addr := range endpoints {
		addr := addr
		go func() {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}
			logger.Bg().Info(
				fmt.Sprintf("debug listen on \"%s\"", listener.Addr()),
				zap.String("addr", addr),
			)
			err = http.Serve(listener, pprofServeMux)
			logger.Bg().Error(
				"http.ListenAndServe(\"%s\", pprofServeMux) error(%v)",
				zap.String("addr", addr),
				zap.Error(err),
			)
			panic(err)
		}()
	}
}
