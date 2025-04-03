package main

import (
	"github.com/clambin/homewizard-exporter/internal/collector"
)

func main() {
	if err := collector.Run(); err != nil {
		panic(err)
	}
}
