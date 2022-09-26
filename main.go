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

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ceph/go-ceph/rgw/admin"
	"github.com/koor-tech/extended-ceph-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	flags                    = flag.NewFlagSet("exporter", flag.ExitOnError)
	defaultEnabledCollectors = "rgw_user_quota,rgw_buckets"
	log                      = logrus.New()
)

type CmdLineOpts struct {
	Version  bool
	LogLevel string

	CollectorsEnabled string

	RGWHost     string
	RGWUser     string
	RGWPassword string

	ListenHost  string
	MetricsPath string

	CachingEnabled bool
	CacheDuration  time.Duration
}

var opts CmdLineOpts

func init() {
	flags.BoolVar(&opts.Version, "version", false, "Show version info and exit")
	flags.StringVar(&opts.LogLevel, "log-level", "INFO", "Set log level")

	flags.StringVar(&opts.CollectorsEnabled, "collectors-enabled", defaultEnabledCollectors, "List of enabled collectors")
	flags.StringVar(&opts.RGWHost, "rgw-host", "", "RGW Host URL")
	flags.StringVar(&opts.RGWUser, "rgw-user", "", "RGW Username")
	flags.StringVar(&opts.RGWPassword, "rgw-password", "", "RGW Password")

	flags.StringVar(&opts.ListenHost, "listen-host", ":9138", "Exporter listen host")
	flags.StringVar(&opts.MetricsPath, "metrics-path", "/metrics", "Set the metrics endpoint path")

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

	if opts.Version {
		fmt.Fprintln(os.Stdout, version.Print(os.Args[0]))
		os.Exit(0)
	}

	log.Out = os.Stdout

	// Set log level
	l, err := logrus.ParseLevel(opts.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(l)

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

	if err = prometheus.Register(NewExtendedCephMetricsCollector(log, collectors, opts.CachingEnabled, opts.CacheDuration)); err != nil {
		log.Fatalf("Couldn't register collector: %s", err)
	}

	log.Infof("Listening on %s", opts.ListenHost)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
<html>
	<head><title>Extended Ceph Exporter</title></head>
	<body>
		<h1>Extended Ceph Exporter</h1>
		<p><a href="` + opts.MetricsPath + `">Metrics</a></p>
	</body>
</html>`))
	})

	handler := promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog:      log,
			ErrorHandling: promhttp.ContinueOnError,
		})

	http.HandleFunc(opts.MetricsPath, func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

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
