package collector

import (
	"codeberg.org/clambin/go-homewizard"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

var (
	currentPower = prometheus.NewDesc(
		prometheus.BuildFQName("homewizard", "current", "power"),
		"Current power (in W)",
		nil, nil,
	)
	currentCurrent = prometheus.NewDesc(
		prometheus.BuildFQName("homewizard", "current", "current"),
		"Current current (in A)",
		nil, nil,
	)
	currentVoltage = prometheus.NewDesc(
		prometheus.BuildFQName("homewizard", "current", "voltage"),
		"Current voltage (in V)",
		nil, nil,
	)
	peakPower = prometheus.NewDesc(
		prometheus.BuildFQName("homewizard", "peak", "power"),
		"Latest peak power (in W)",
		nil, nil,
	)
)

var _ prometheus.Collector = &Collector{}

type Collector struct {
	Client HomeWizardClient
	Logger *slog.Logger
}

type HomeWizardClient interface {
	GetRecentMeasurement(ctx context.Context) (*homewizard.RecentMeasurement, error)
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- currentPower
	ch <- currentCurrent
	ch <- currentVoltage
	ch <- peakPower
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	measurement, err := c.Client.GetRecentMeasurement(context.Background())
	if err != nil {
		c.Logger.Error("failed to collect homewizard metrics", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(currentPower, prometheus.GaugeValue, measurement.ActivePowerW)
	ch <- prometheus.MustNewConstMetric(currentCurrent, prometheus.GaugeValue, measurement.ActiveCurrentL1A)
	ch <- prometheus.MustNewConstMetric(currentVoltage, prometheus.GaugeValue, measurement.ActiveCurrentL1A)
	ch <- prometheus.MustNewConstMetric(peakPower, prometheus.GaugeValue, measurement.MontlyPowerPeakW)
}
