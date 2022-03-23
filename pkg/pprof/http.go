package pprof

import (
	"github.com/tkeel-io/rule-util/pkg/debug"
	"github.com/tkeel-io/rule-util/pkg/log"
	"github.com/tkeel-io/rule-util/pkg/logfield"
	"net/http"
	"net/http/pprof"
	"sync"
)

// Init pprof
// endpoints 127.0.0.1:7001,10.0.0.1:7001
// endpoints 0.0.0.0:7001

var once sync.Once

func Init(endpoints []string) {
	once.Do(func() {
		pprofServeMux := http.NewServeMux()
		pprofServeMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofServeMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		pprofServeMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofServeMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofServeMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		for _, addr := range endpoints {
			go func() {
				if err := http.ListenAndServe(addr, pprofServeMux); err != nil {
					log.Error("http.ListenAndServe",
						logf.String("addr", addr),
						logf.Error(err))
					panic(err)
				} else {
					log.Info("pprof listen",
						logf.String("addr", addr))
				}
			}()
		}
	})
}

func HandleFunc() []*debug.HandleFunc {
	return []*debug.HandleFunc{
		{"/pprof/", pprof.Index},
		{"/pprof/cmdline", pprof.Cmdline},
		{"/pprof/profile", pprof.Profile},
		{"/pprof/symbol", pprof.Symbol},
	}
}
