package v1

import (
	"net/http"
	"net/http/pprof"
	"runtime"
)

// newProfilerResource injects endpoints of net/http/pprof into the provided router.
func newProfilerResource(router *http.ServeMux) {
	// NB (alkurbatov): Tweak memory profiling rate to spot more allocations.
	runtime.MemProfileRate = 2048

	router.Handle("GET /", http.RedirectHandler("/debug/pprof/", http.StatusMovedPermanently))
	router.Handle("GET /pprof", http.RedirectHandler("/debug/pprof/", http.StatusMovedPermanently))

	router.HandleFunc("GET /debug/pprof/", pprof.Index)
	router.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	router.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
}
