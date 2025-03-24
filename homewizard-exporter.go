package main

import (
	"codeberg.org/clambin/go-homewizard"
	"errors"
	"flag"
	"github.com/clambin/homewizard-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var (
	version = "change-me"

	addr   = flag.String("addr", ":8080", "Prometheus exporter address")
	target = flag.String("target", "http://192.168.0.188", "URL of the homewizard meter")
	debug  = flag.Bool("debug", false, "Enable debug logging")
)

func main() {
	if !strings.HasSuffix(*target, "/") {
		*target += "/"
	}

	var opt slog.HandlerOptions
	if *debug {
		opt.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &opt))

	logger.Info("Starting homewizard exporter", "version", version)

	c := collector.Collector{
		Client: &homewizard.Client{
			HTTPClient: http.DefaultClient,
			Target:     *target,
		},
		Logger: logger,
	}
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(*addr, nil); !errors.Is(err, http.ErrServerClosed) {
		logger.Error("error serving metrics", "err", err)
	}
}
