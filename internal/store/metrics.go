package store

import "github.com/prometheus/client_golang/prometheus"

var (
	getTotalOps = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kvstore_get_total",
		Help: "Total number of GET operations",
	})

	setTotalOps = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kvstore_set_total",
		Help: "Total number of SET operations",
	})

	deleteTotalOps = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kvstore_delete_total",
		Help: "Total number of DELETE operations",
	})

	expiredKeysTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kvstore_expired_keys_total",
		Help: "Total number of keys deleted due to TTL expiry",
	})

	activeKeys = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "kvstore_active_keys",
		Help: "Current number of keys in the store",
	})
)

func init() {
	prometheus.MustRegister(getTotalOps)
	prometheus.MustRegister(setTotalOps)
	prometheus.MustRegister(deleteTotalOps)
	prometheus.MustRegister(expiredKeysTotal)
	prometheus.MustRegister(activeKeys)
}
