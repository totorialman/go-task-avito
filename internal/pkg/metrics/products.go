package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ProductMetrics struct {
	ReTotal      prometheus.Counter
	PVZTotal     prometheus.Counter
	ProductTotal prometheus.Counter
}

func NewProductMetrics() (*ProductMetrics, error) {
	var metr ProductMetrics
	metr.ReTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ReTotal",
			Help: "Number of total hits.",
		})
	if err := prometheus.Register(metr.ReTotal); err != nil {
		return nil, err
	}
	metr.PVZTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "PVZTotal",
			Help: "Number of total hits.",
		})
	if err := prometheus.Register(metr.PVZTotal); err != nil {
		return nil, err
	}
	metr.ProductTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ProductTotal",
			Help: "Number of total hits.",
		})
	if err := prometheus.Register(metr.ProductTotal); err != nil {
		return nil, err
	}

	return &metr, nil
}
func (m *ProductMetrics) IncreaseHitsReTotal() {
	m.ReTotal.Inc()
}

func (m *ProductMetrics) IncreaseHitsPVZTotal() {
	m.PVZTotal.Inc()
}

func (m *ProductMetrics) IncreaseHitsProductTotal() {
	m.ProductTotal.Inc()
}
