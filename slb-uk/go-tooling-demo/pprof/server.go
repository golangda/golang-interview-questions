package main

import (
	_ "net/http/pprof"
	"fmt"
	"log"
	"math"
	"net/http"
	"runtime"
	"time"
)

// burnCPU spins on some floating point work for a duration to simulate load.
func burnCPU(d time.Duration) {
	end := time.Now().Add(d)
	x := 0.0
	for time.Now().Before(end) {
		// some work; value discarded to avoid dead-code elimination
		x += math.Sin(float64(time.Now().UnixNano()))
	}
	_ = x
}

func workHandler(w http.ResponseWriter, r *http.Request) {
	// optional duration query: /work?sec=2
	sec := 2 * time.Second
	if v := r.URL.Query().Get("sec"); v != "" {
		if n, err := time.ParseDuration(v + "s"); err == nil {
			sec = n
		}
	}
	burnCPU(sec)
	fmt.Fprintln(w, "done in", sec)
}

func main() {
	// Enable additional profiles
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	// Register a simple workload handler on the default mux (same mux pprof uses)
	http.HandleFunc("/work", workHandler)

	log.Println("Serving pprof + demo at http://localhost:6060")
	log.Println("Try: curl http://localhost:6060/work")
	log.Println("CPU profile: go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30")

	// Start HTTP server with pprof endpoints on :6060 using the default mux
	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
		log.Fatal(err)
	}
}
