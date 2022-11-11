package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {

	flag.Parse()
	reg := prometheus.NewRegistry()

	opsQueued := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "sdewan_system",
		Subsystem:   "pkcs11_hsm",
		Name:        "key_pair_create",
		ConstLabels: prometheus.Labels{"type": "sgx"},
		Help:        "Number of key pair to be created",
	}, []string{"operation"})
	err := reg.Register(opsQueued)
	if err != nil {
		log.Fatal("register failed", err)

	}
	go func() {
		var i int
		for {
			if i%2 == 0 {
				opsQueued.WithLabelValues("create").Inc()
			} else {
				opsQueued.WithLabelValues("delete").Inc()
			}
			i++
			time.Sleep(time.Second * 1)
		}

	}()

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(*addr, nil))
	return
}
