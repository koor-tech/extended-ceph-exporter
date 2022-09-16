package collector

import (
	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/prometheus/client_golang/prometheus"
)

const Namespace = "ceph"

type Clients struct {
	RGWAdminAPI *admin.API
}

type Collector interface {
	Update(chan<- prometheus.Metric) error
}

type NewCollectorFunc func(*Clients) (Collector, error)

var Factories map[string]NewCollectorFunc = map[string]NewCollectorFunc{}
