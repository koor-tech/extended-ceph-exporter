/*
Copyright 2022 Koor Technologies, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"context"

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
		return err
	}

	for _, bucketName := range buckets {
		bucketInfo, err := c.api.GetBucketInfo(context.Background(), admin.Bucket{
			Bucket: bucketName,
		})
		if err != nil {
			return err
		}

		labels := map[string]string{
			"bucket": bucketName,
			"uid":    bucketInfo.Owner,
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size"),
			"RGW Bucket Size",
			nil, labels)
		if bucketInfo.Usage.RgwMain.Size == nil {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, 0.0)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.Size))
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size_kb"),
			"RGW Bucket Size actual",
			nil, labels)
		if bucketInfo.Usage.RgwMain.SizeKb == nil {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, 0.0)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.SizeKb))
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size_kb_actual"),
			"RGW Bucket Size KiB actual",
			nil, labels)
		if bucketInfo.Usage.RgwMain.SizeKbActual == nil {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, 0.0)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.SizeKbActual))
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_size_kb_utilized"),
			"RGW Bucket Size KiB utilized",
			nil, labels)
		if bucketInfo.Usage.RgwMain.SizeKbUtilized == nil {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, 0.0)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.SizeKbUtilized))
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_num_objects"),
			"RGW Bucket Num Objects",
			nil, labels)
		if bucketInfo.Usage.RgwMain.NumObjects == nil {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, 0.0)
		} else {
			ch <- prometheus.MustNewConstMetric(
				c.current, prometheus.GaugeValue, float64(*bucketInfo.Usage.RgwMain.NumObjects))
		}

		if bucketInfo.BucketQuota.Enabled == nil || !*bucketInfo.BucketQuota.Enabled {
			continue
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_quota_max_size_kb"),
			"RGW Bucket Quota Max Size KiB",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*bucketInfo.BucketQuota.MaxSizeKb))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "bucket_quota_max_objects"),
			"RGW Bucket Quota Max Objects",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*bucketInfo.BucketQuota.MaxObjects))

	}

	return nil
}
