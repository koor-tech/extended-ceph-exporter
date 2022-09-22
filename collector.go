package main

import (
	"sync"
	"time"

	"github.com/koor-tech/extended-cephmetrics-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "collector_duration_seconds"),
		"Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "scrape", "collector_success"),
		"Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

// ExtendedCephMetricsCollector contains the collectors to be used
type ExtendedCephMetricsCollector struct {
	log             *logrus.Logger
	lastCollectTime time.Time
	collectors      map[string]collector.Collector

	// Cache related
	cachingEnabled bool
	cacheDuration  time.Duration
	cache          []prometheus.Metric
	cacheMutex     sync.Mutex
}

func NewExtendedCephMetricsCollector(log *logrus.Logger, collectors map[string]collector.Collector, cachingEnabled bool, cacheDuration time.Duration) *ExtendedCephMetricsCollector {
	return &ExtendedCephMetricsCollector{
		log:             log,
		cache:           make([]prometheus.Metric, 0),
		lastCollectTime: time.Unix(0, 0),
		collectors:      collectors,
		cachingEnabled:  cachingEnabled,
		cacheDuration:   cacheDuration,
	}
}

// Describe implements the prometheus.Collector interface.
func (n *ExtendedCephMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect implements the prometheus.Collector interface.
func (n *ExtendedCephMetricsCollector) Collect(outgoingCh chan<- prometheus.Metric) {
	if n.cachingEnabled {
		n.cacheMutex.Lock()
		defer n.cacheMutex.Unlock()

		expiry := n.lastCollectTime.Add(n.cacheDuration)
		if time.Now().Before(expiry) {
			n.log.Debugf("Using cache. Now: %s, Expiry: %s, LastCollect: %s", time.Now().String(), expiry.String(), n.lastCollectTime.String())
			for _, cachedMetric := range n.cache {
				n.log.Debugf("Pushing cached metric %s to outgoingCh", cachedMetric.Desc().String())
				outgoingCh <- cachedMetric
			}
			return
		}
		// Clear cache, but keep slice
		n.cache = n.cache[:0]
	}

	metricsCh := make(chan prometheus.Metric)

	// Wait to ensure outgoingCh is not closed before the goroutine is finished
	wgOutgoing := sync.WaitGroup{}
	wgOutgoing.Add(1)
	go func() {
		for metric := range metricsCh {
			outgoingCh <- metric
			if n.cachingEnabled {
				n.log.Debugf("Appending metric %s to cache", metric.Desc().String())
				n.cache = append(n.cache, metric)
			}
		}
		n.log.Debug("Finished pushing metrics from metricsCh to outgoingCh")
		wgOutgoing.Done()
	}()

	wgCollection := sync.WaitGroup{}
	wgCollection.Add(len(n.collectors))
	for name, coll := range n.collectors {
		go func(name string, coll collector.Collector) {
			begin := time.Now()
			err := coll.Update(metricsCh)
			duration := time.Since(begin)
			var success float64

			if err != nil {
				n.log.Errorf("%s collector failed after %fs: %s", name, duration.Seconds(), err)
				success = 0
			} else {
				n.log.Debugf("%s collector succeeded after %fs.", name, duration.Seconds())
				success = 1
			}
			metricsCh <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
			metricsCh <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
			wgCollection.Done()
		}(name, coll)
	}

	n.log.Debug("Waiting for collectors")
	wgCollection.Wait()
	n.log.Debug("Finished waiting for collectors")

	n.lastCollectTime = time.Now()
	n.log.Debugf("Updated lastCollectTime to %s", n.lastCollectTime.String())

	close(metricsCh)

	n.log.Debug("Waiting for outgoing Adapter")
	wgOutgoing.Wait()
	n.log.Debug("Finished waiting for outgoing Adapter")
}
