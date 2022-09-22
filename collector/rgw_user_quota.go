package collector

import (
	"context"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	Factories["rgw_user_quota"] = NewRGWUserQuota
}

func NewRGWUserQuota(c *Clients) (Collector, error) {
	return &RGWUserQuota{
		api: c.RGWAdminAPI,
	}, nil
}

type RGWUserQuota struct {
	api *admin.API

	current *prometheus.Desc
}

func (c *RGWUserQuota) Update(ch chan<- prometheus.Metric) error {
	// Get the "admin" user
	users, err := c.api.GetUsers(context.Background())
	if err != nil {
		return err
	}

	// Iterate over users to get quota
	for _, user := range *users {
		userQuota, err := c.api.GetUserQuota(context.Background(), admin.QuotaSpec{
			UID: user,
		})
		if err != nil {
			return err
		}

		if userQuota.Enabled == nil {
			continue
		}

		labels := map[string]string{
			"uid": user,
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "user_userQuota_max_size"),
			"RGW User Quota max size",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*userQuota.MaxSize))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "user_quota_max_size_kb"),
			"RGW User Quota max size KiB",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*userQuota.MaxSizeKb))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "user_quota_max_objects"),
			"RGW User Quota max objects",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*userQuota.MaxObjects))
	}

	return nil
}
