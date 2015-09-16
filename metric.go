package main

import (
	"os"
	"time"

	"github.com/zorkian/go-datadog-api"
)

var hostname, _ = os.Hostname()
var now = float64(time.Now().Unix())

func metricInt(name string, val int, tags []string) datadog.Metric {
	return metric(name, float64(val), tags)
}

func metricInt32(name string, val int32, tags []string) datadog.Metric {
	return metric(name, float64(val), tags)
}

func metricInt64(name string, val int64, tags []string) datadog.Metric {
	return metric(name, float64(val), tags)
}

func metricUint(name string, val uint, tags []string) datadog.Metric {
	return metric(name, float64(val), tags)
}

func metricUint64(name string, val uint64, tags []string) datadog.Metric {
	return metric(name, float64(val), tags)
}

func metric(name string, val float64, tags []string) datadog.Metric {
	return datadog.Metric{
		Metric: "nsqdog." + name,
		Points: []datadog.DataPoint{{now, val}},
		Type:   "gauge",
		Host:   hostname,
		Tags:   tags,
	}
}
