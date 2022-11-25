package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"metric-exporter-demo/telemetry"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var (
	keyOpsMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "sdewan_system",
		Subsystem:   "pkcs11_hsm",
		Name:        "key_pair_create",
		ConstLabels: prometheus.Labels{"type": "sgx"},
		Help:        "Number of key pair to be created",
	}, []string{"operation"})

	healthyCheckMetric = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   "sdewan_system",
		Subsystem:   "pkcs11_hsm",
		Name:        "health_check",
		ConstLabels: prometheus.Labels{"type": "sgx"},
		Help:        "Number of health check",
	})
)

func main() {

	flag.Parse()
	reg := prometheus.NewRegistry()

	_ = reg.Register(keyOpsMetric)
	_ = reg.Register(healthyCheckMetric)

	_, err := telemetry.InitProvider()
	if err != nil {
		log.Println(err)
	}

	// Mock a QPS
	go func() {
		var i int
		for {
			keyOpsMetric.WithLabelValues("create").Inc()
			keyOpsMetric.WithLabelValues("delete").Inc()
			i++
			time.Sleep(time.Second * 30)
		}
	}()
	log.Println("Service starting...")

	// Expose /metrics HTTP endpoint using the created custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.Handle("/go-metrics", promhttp.Handler())
	http.HandleFunc("/ping", handlerPing)
	http.HandleFunc("/job", handleJob)
	log.Fatal(http.ListenAndServe(*addr, nil))
	return
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		return
	}
	healthyCheckMetric.Inc()
	log.Println(time.Now(), r.Method, r.RequestURI, r.UserAgent(), "service healthy check!")
}

func handleJob(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("success"))
	if err != nil {
		return
	}
	ctx := context.TODO()
	// Use the global TracerProvider.
	tr := otel.Tracer("fibonacci-job")
	// work begin
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("bar-key").String("bar-value"))
	defer span.End()
	for i := 0; i < 10; i++ {
		c, iSpan := tr.Start(ctx, fmt.Sprintf("Sample-%d", i))
		log.Printf("Doing really hard work (%d / 10)\n", i+1)
		<-time.After(time.Millisecond * 200)
		iSpan.End()
		ctx = c
	}
}
