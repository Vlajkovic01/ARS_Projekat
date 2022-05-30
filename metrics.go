package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	currentCount = 0

	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_http_hit_total",
			Help: "Total number of http hits.",
		},
	)

	postConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_config_hit_total",
			Help: "Total number of create config hits.",
		},
	)

	postConfigVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_config_ver_hit_total",
			Help: "Total number of add new config version hits.",
		},
	)

	getConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_config_hit_total",
			Help: "Total number of get one config hits.",
		},
	)

	getConfigVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_config_ver_hit_total",
			Help: "Total number of get all config versions hits.",
		},
	)

	postGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_group_hit_total",
			Help: "Total number of post new group hits.",
		},
	)

	postGroupVersionHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_post_group_ver_hit_total",
			Help: "Total number of adding new group version hits.",
		},
	)

	getGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_group_hit_total",
			Help: "Total number of get group hits.",
		},
	)

	getGroupConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_get_group_config_hit_total",
			Help: "Total number of get group configs by label hits.",
		},
	)

	addGroupConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_add_group_config_hit_total",
			Help: "Total number of add new config to a group hits.",
		},
	)

	deleteGroupHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_del_group_hit_total",
			Help: "Total number of delete group hits.",
		},
	)

	deleteConfigHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configstore_del_config_hit_total",
			Help: "Total number of delete config hits.",
		},
	)

	metricsList = []prometheus.Collector{
		postConfigHits, getConfigVersionHits, postConfigVersionHits, getConfigHits,
		deleteConfigHits, postGroupHits, postGroupVersionHits, getGroupHits, deleteGroupHits,
		getGroupConfigHits, addGroupConfigHits, httpHits,
	}

	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	prometheusRegistry.MustRegister(metricsList...)
}

func metricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func countPostConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postConfigHits.Inc()
		f(w, r) // original function call
	}
}

func countGetConfigVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigVersionHits.Inc()
		f(w, r) // original function call
	}
}

func countPostConfigVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postConfigVersionHits.Inc()
		f(w, r) // original function call
	}
}

func countGetConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getConfigHits.Inc()
		f(w, r) // original function call
	}
}

func countDeleteConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		deleteConfigHits.Inc()
		f(w, r) // original function call
	}
}

func countPostGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postGroupHits.Inc()
		f(w, r) // original function call
	}
}

func countPostGroupVersion(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		postGroupVersionHits.Inc()
		f(w, r) // original function call
	}
}

func countGetGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupHits.Inc()
		f(w, r) // original function call
	}
}

func countDeleteGroup(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		deleteGroupHits.Inc()
		f(w, r) // original function call
	}
}

func countGetGroupConfigs(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		getGroupConfigHits.Inc()
		f(w, r) // original function call
	}
}

func countAddGroupConfig(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		httpHits.Inc()
		addGroupConfigHits.Inc()
		f(w, r) // original function call
	}
}
