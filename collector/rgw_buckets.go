package collector

import (
	"context"
	"fmt"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	Factories["rgw_buckets"] = NewRGWBuckets
}

func NewRGWBuckets(c *Clients) (Collector, error) {
	return &RGWBuckets{
		api: c.RGWAdminAPI,
	}, nil
}

type RGWBuckets struct {
	api *admin.API

	current *prometheus.Desc
}

func (c *RGWBuckets) Update(ch chan<- prometheus.Metric) error {
	buckets, err := c.api.ListBuckets(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("buckets %s \n ", buckets)

	for _, bucketName := range buckets {
		bucketInfo, err := c.api.GetBucketInfo(context.Background(), admin.Bucket{
			Bucket: bucketName,
		})
		if err != nil {
			return err
		}

		labels := map[string]string{
			"bucket": bucketName,
			"owner":  bucketInfo.Owner,
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size"),
			"RGW Bucket Size",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.Size))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size_kb"),
			"RGW Bucket Size actual",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.SizeKb))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size_kb_actual"),
			"RGW Bucket Size kb actual",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.SizeKbActual))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size_kb_utilized"),
			"RGW Bucket Size kb utilized",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.SizeKbUtilized))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_num_objects"),
			"RGW Bucket Num Objects",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.NumObjects))
	}

	return nil
}
