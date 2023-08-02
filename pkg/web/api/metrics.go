package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"meteor-server/pkg/db"
	"net/http"
)

var handler http.Handler

func InitMetrics() {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "meteor_downloads_total",
			Help: "Total number of downloads",
		}, func() float64 {
			return float64(db.GetGlobal().Downloads)
		}),
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "meteor_online_players_total",
			Help: "Total number of online players",
		}, func() float64 {
			return float64(GetPlayingCount())
		}),
		promauto.NewCounterFunc(prometheus.CounterOpts{
			Name: "meteor_donators_total",
			Help: "Total number of accounts with donator status",
		}, func() float64 {
			return float64(db.GetDonatorCount())
		}),
	)

	handler = promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	handler.ServeHTTP(w, r)
}
