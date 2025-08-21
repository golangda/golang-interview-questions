package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Demo app metrics
var (
	appRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Total number of demo endpoint requests",
		},
		[]string{"path"},
	)

	appWorkDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "app_work_duration_seconds",
			Help:    "Simulated work duration",
			Buckets: prometheus.DefBuckets,
		},
	)

	appInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_inflight_requests",
			Help: "In-flight requests for demo endpoints",
		},
	)
)

func main() {
	// Register standard collectors explicitly (process_*, go_*, build info)
	// Default registry already has Go + process collectors.
	prometheus.MustRegister(collectors.NewBuildInfoCollector())


	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("/work", withMetrics("/work", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() { appWorkDuration.Observe(time.Since(start).Seconds()) }()

		n := 20000 + rand.Intn(20000)
		sum := 0
		for i := 0; i < n; i++ {
			sum += i * (i % 7)
		}
		time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
		fmt.Fprintf(w, "did some work, sum=%d\n", sum)
	}))

	mux.HandleFunc("/alloc", withMetrics("/alloc", func(w http.ResponseWriter, r *http.Request) {
		bufs := make([][]byte, 100)
		for i := range bufs {
			bufs[i] = make([]byte, 1<<16) // 64 KiB
			for j := range bufs[i] {
				bufs[i][j] = byte(j)
			}
		}
		fmt.Fprintf(w, "allocated ~%d KiB then released\n", len(bufs)*64)
	}))

	mux.HandleFunc("/goroutines", withMetrics("/goroutines", func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		n := 100 + rand.Intn(200)
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				defer wg.Done()
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			}()
		}
		fmt.Fprintf(w, "spawned %d goroutines; currently: %d\n", n, runtime.NumGoroutine())
		wg.Wait()
	}))

	// Instrumented /metrics to expose promhttp_* metrics
	metricsHandler := promhttp.Handler()
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		metricsHandler,
	))

	addr := ":2112"
	log.Printf("Prometheus demo listening on %s", addr)
	log.Printf("Try: http://localhost%[1]s/metrics, /work, /alloc, /goroutines, /healthz", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func withMetrics(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appInFlight.Inc()
		defer appInFlight.Dec()
		appRequestsTotal.WithLabelValues(path).Inc()
		next.ServeHTTP(w, r)
	}
}