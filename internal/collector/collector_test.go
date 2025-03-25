package collector

import (
	"codeberg.org/clambin/go-homewizard"
	"context"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"log/slog"
	"strings"
	"testing"
)

func TestCollector(t *testing.T) {
	c := Collector{
		Client: &fakeClient{},
		Logger: slog.New(slog.DiscardHandler),
	}

	if err := testutil.CollectAndCompare(&c, strings.NewReader(`
# HELP homewizard_current_current Current current (in A)
# TYPE homewizard_current_current gauge
homewizard_current_current 10

# HELP homewizard_current_power Current power (in W)
# TYPE homewizard_current_power gauge
homewizard_current_power 2400

# HELP homewizard_current_voltage Current voltage (in V)
# TYPE homewizard_current_voltage gauge
homewizard_current_voltage 240

# HELP homewizard_peak_power Latest peak power (in W)
# TYPE homewizard_peak_power gauge
homewizard_peak_power 6000
`)); err != nil {
		t.Errorf("unexpected collecting result:\n%v", err)
	}
}

var _ HomeWizardClient = &fakeClient{}

type fakeClient struct{}

func (f fakeClient) GetRecentMeasurement(_ context.Context) (homewizard.RecentMeasurement, error) {
	return homewizard.RecentMeasurement{MontlyPowerPeakW: 6000, ActiveVoltageL1V: 240, ActiveCurrentL1A: 10, ActivePowerW: 2400}, nil
}

func (f fakeClient) GetDeviceInformation(_ context.Context) (homewizard.DeviceInformation, error) {
	return homewizard.DeviceInformation{}, nil
}
