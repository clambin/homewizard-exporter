package collector

import (
	"codeberg.org/clambin/go-homewizard"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"sync"
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
	Once   sync.Once
}

type HomeWizardClient interface {
	GetRecentMeasurement(ctx context.Context) (homewizard.RecentMeasurement, error)
	GetDeviceInformation(ctx context.Context) (homewizard.DeviceInformation, error)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- currentPower
	ch <- currentCurrent
	ch <- currentVoltage
	ch <- peakPower
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.Once.Do(c.logDeviceInfo)

	measurement, err := c.Client.GetRecentMeasurement(context.Background())
	if err != nil {
		c.Logger.Error("failed to collect homewizard metrics", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(currentPower, prometheus.GaugeValue, measurement.ActivePowerW)
	ch <- prometheus.MustNewConstMetric(currentCurrent, prometheus.GaugeValue, measurement.ActiveCurrentL1A)
	ch <- prometheus.MustNewConstMetric(currentVoltage, prometheus.GaugeValue, measurement.ActiveVoltageL1V)
	ch <- prometheus.MustNewConstMetric(peakPower, prometheus.GaugeValue, measurement.MontlyPowerPeakW)
}

func (c *Collector) logDeviceInfo() {
	info, err := c.Client.GetDeviceInformation(context.Background())
	if err != nil {
		c.Logger.Error("failed to collect homewizard device information", "err", err)
		return
	}
	c.Logger.Info("Device found", slog.Group("device",
		"name", info.ProductName,
		"type", info.ProductType,
		"firmware", info.FirmwareVersion,
		"api", info.ApiVersion,
	))
}
