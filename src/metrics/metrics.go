package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var (
	HueCacheUpdateLatencyMSHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "hue_cache_update_latency_ms",
		Help:    "The length of time (in milliseconds) it takes to update the HUE cache",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 5000},
	})

	HueCacheUpdateErrorsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hue_cache_update_errors",
			Help: "The number of times attempting to update a HUE cache fails",
		},
		[]string{"type"},
	)

	HueDeviceStateChangeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "hue_device_state_change",
			Help: "The current value of a HUE device state",
		},
		[]string{"name", "type", "state", "sensorType"},
	)

	HueEventsEmittedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hue_events_emitted",
			Help: "The number of HUE events emitted",
		},
		[]string{"event", "type"},
	)
)

func Startup(port int, log *logrus.Logger) {
	http.Handle("/metrics", promhttp.Handler())

	prometheus.MustRegister(HueCacheUpdateLatencyMSHistogram)
	prometheus.MustRegister(HueCacheUpdateErrorsCounter)
	prometheus.MustRegister(HueDeviceStateChangeGauge)
	prometheus.MustRegister(HueEventsEmittedCounter)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
