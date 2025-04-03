package collector

import (
	"codeberg.org/clambin/go-homewizard"
	"errors"
	"flag"
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
	target = flag.String("target", "", "URL of the homewizard meter")
	debug  = flag.Bool("debug", false, "Enable debug logging")
)

func Run() error {
	flag.Parse()
	if *target == "" {
		return errors.New("target flag is required")
	}
	if !strings.HasSuffix(*target, "/") {
		*target += "/"
	}
	client := &homewizard.Client{
		HTTPClient: http.DefaultClient,
		Target:     *target,
	}
	var opt slog.HandlerOptions
	if *debug {
		opt.Level = slog.LevelDebug
	}
	return run(client, prometheus.DefaultRegisterer, *addr, slog.New(slog.NewTextHandler(os.Stderr, &opt)))
}

func run(client HomeWizardClient, registry prometheus.Registerer, addr string, logger *slog.Logger) error {
	logger.Info("Starting homewizard exporter", "version", version)

	c := Collector{Client: client, Logger: logger}
	registry.MustRegister(&c)

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}
