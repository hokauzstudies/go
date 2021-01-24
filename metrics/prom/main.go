package prom

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type service struct {
	counter   *prometheus.CounterVec
	summary   *prometheus.SummaryVec
	histogram *prometheus.HistogramVec
}

func new(coll collector) {
	s := service{}

	for _, metric := range coll.Metrics {
		s.append(metric)
	}
}

func (s *service) append(metric metrictype) {
	if metric.Type == "counter" {
		s.counter = makeCounter(metric)
		register(metric.Name, s.counter)
		return
	}

	if metric.Type == "summary" {
		s.summary = makeSummary(metric)
		register(metric.Name, s.summary)
		return
	}

	if metric.Type == "histogram" {
		s.histogram = makeHistogram(metric)
		register(metric.Name, s.histogram)
		return
	}
}

func makeCounter(metric metrictype) *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "requests",  // TODO acept namespae
		Name:      metric.Name, // TODO concacternar com namespace
		Help:      metric.Helper,
	}, metric.Labels)
}

func makeSummary(metric metrictype) *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "requests",  // TODO acept namespae
		Name:      metric.Name, // TODO concacternar com namespace
		Help:      metric.Helper,
	}, metric.Labels)
}

func makeHistogram(metric metrictype) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "requests",  // TODO acept namespae
		Name:      metric.Name, // TODO concacternar com namespace
		Help:      metric.Helper,
		Buckets:   prometheus.DefBuckets,
	}, metric.Labels)
}

func register(name string, coll prometheus.Collector) {
	err := prometheus.Register(coll)
	if err != nil {
		fmt.Println("error", name, err)
	}
}

func (s *service) trigger() {
	s.counter.WithLabelValues("code", "method").Inc()
}
