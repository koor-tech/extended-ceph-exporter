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
		panic(err)
	}

	// Iterate over users to get quota
	for _, user := range *users {
		quota, err := c.api.GetUserQuota(context.Background(), admin.QuotaSpec{
			UID: user,
		})
		if err != nil {
			return err
		}

		labels := map[string]string{
			"uid": user,
		}

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "user_quota_max_size_kb"),
			"RGW User Quota max size kb",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*quota.MaxSizeKb))

		c.current = prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rgw", "user_quota_max_objects"),
			"RGW User Quota max objects",
			nil, labels)
		ch <- prometheus.MustNewConstMetric(
			c.current, prometheus.GaugeValue, float64(*quota.MaxObjects))
	}

	return nil
}
