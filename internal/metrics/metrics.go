package metrics

import "github.com/prometheus/client_golang/prometheus"

var IncomingTraffic = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "incoming_traffic",
	Help: "Incoming traffic to the application",
})
