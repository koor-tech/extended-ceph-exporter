package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/koor-tech/extended-cephmetrics-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var flags = flag.NewFlagSet("exporter", flag.ExitOnError)

var defaultEnabledCollectors = "rgw_user_quota"

type CmdLineOpts struct {
	CollectorsEnabled string

	RGWHost     string
	RGWUser     string
	RGWPassword string

	ListenHost string

	CachingEnabled bool
	CacheDuration  time.Duration
}

var opts CmdLineOpts

func init() {
	flags.StringVar(&opts.CollectorsEnabled, "collectors-enabled", defaultEnabledCollectors, "List of enabled collectors")
	flags.StringVar(&opts.RGWHost, "rgw-host", "", "RGW Host URL")
	flags.StringVar(&opts.RGWUser, "rgw-user", "", "RGW Username")
	flags.StringVar(&opts.RGWPassword, "rgw-password", "", "RGW Password")
	flags.StringVar(&opts.ListenHost, "listen-host", ":9138", "Exporter listen host")

	flags.BoolVar(&opts.CachingEnabled, "cache-enabled", false, "Enable metrics caching to reduce load")
	flags.DurationVar(&opts.CacheDuration, "cache-duration", 20*time.Second, "Cache duration in seconds")
}

func flagNameFromEnvName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func parseFlagsAndEnvVars() error {
	for _, v := range os.Environ() {
		vals := strings.SplitN(v, "=", 2)

		if !strings.HasPrefix(vals[0], "CEPH_METRICS_") {
			continue
		}
		flagName := flagNameFromEnvName(strings.ReplaceAll(vals[0], "CEPH_METRICS_", ""))

		fn := flags.Lookup(flagName)
		if fn == nil || fn.Changed {
			continue
		}

		if err := fn.Value.Set(vals[1]); err != nil {
			return err
		}
		fn.Changed = true
	}

	return flags.Parse(os.Args[1:])
}

func main() {
	if err := parseFlagsAndEnvVars(); err != nil {
		log.Fatal(err)
	}

	// Check if required flags are given
	for _, k := range []string{"rgw-host", "rgw-user", "rgw-password"} {
		flag := flags.Lookup(k)
		if flag == nil {
			log.Fatalf("flag %s not found during lookup", k)
		}
		if !flag.Changed {
			log.Fatalf("required flag %s not set", flag.Name)
		}
	}

	rgwAdminAPI, err := CreateRGWAPIConnection()
	if err != nil {
		log.Fatal(err)
	}

	clients := &collector.Clients{
		RGWAdminAPI: rgwAdminAPI,
	}

	collectors, err := loadCollectors(opts.CollectorsEnabled, clients)
	if err != nil {
		log.Fatalf("Couldn't load collectors: %s", err)
	}
	log.Infof("Enabled collectors:")
	for n := range collectors {
		log.Infof(" - %s", n)
	}

	if err = prometheus.Register(NewExtendedCephMetricsCollector(collectors, opts.CachingEnabled, opts.CacheDuration)); err != nil {
		log.Fatalf("Couldn't register collector: %s", err)
	}

	log.Infof("Listening on %s", opts.ListenHost)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(opts.ListenHost, nil)
}

func CreateRGWAPIConnection() (*admin.API, error) {
	// Generate a connection object
	co, err := admin.New(opts.RGWHost, opts.RGWUser, opts.RGWPassword, nil)
	if err != nil {
		return nil, err
	}

	return co, nil
}

func loadCollectors(list string, clients *collector.Clients) (map[string]collector.Collector, error) {
	collectors := map[string]collector.Collector{}
	for _, name := range strings.Split(list, ",") {
		fn, ok := collector.Factories[name]
		if !ok {
			return nil, fmt.Errorf("collector '%s' not available", name)
		}
		c, err := fn(clients)
		if err != nil {
			return nil, err
		}
		collectors[name] = c
	}
	return collectors, nil
}
