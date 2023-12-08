// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type LatencyMetric struct {
	Min    time.Duration `json:"Min"`
	Avg    time.Duration `json:"Avg"`
	Max    time.Duration `json:"Max"`
	Perc50 time.Duration `json:"Perc50"`
	Perc90 time.Duration `json:"Perc90"`
	Perc99 time.Duration `json:"Perc99"`
}

func (metric *LatencyMetric) ToPerfData(labels map[string]string) DataItem {
	return DataItem{
		Data: map[string]float64{
			// Let's only export percentiles
			// Max is skewing results and doesn't make much sense to keep track of
			"Perc50": float64(metric.Perc50) / float64(time.Microsecond),
			"Perc90": float64(metric.Perc90) / float64(time.Microsecond),
			"Perc99": float64(metric.Perc99) / float64(time.Microsecond),
		},
		Unit:   "us",
		Labels: labels,
	}
}

type TransactionRateMetric struct {
	TransactionRate float64 `json:"Rate"` // Ops per second
}

func (metric *TransactionRateMetric) ToPerfData(labels map[string]string) DataItem {
	return DataItem{
		Data: map[string]float64{
			"Throughput": metric.TransactionRate,
		},
		Unit:   "ops/s",
		Labels: labels,
	}
}

type ThroughputMetric struct {
	Throughput float64 `json:"Throughput"` // Throughput in bytes/s
}

func (metric *ThroughputMetric) ToPerfData(labels map[string]string) DataItem {
	return DataItem{
		Data: map[string]float64{
			"Throughput": metric.Throughput / 1000000,
		},
		Unit:   "Mb/s",
		Labels: labels,
	}
}

type PerfResult struct {
	Timestamp             time.Time
	Latency               *LatencyMetric
	TransactionRateMetric *TransactionRateMetric
	ThroughputMetric      *ThroughputMetric
}

type PerfTests struct {
	Tool     string
	Test     string
	SameNode bool
	Scenario string
	Sample   int
	Duration time.Duration
}

type PerfSummary struct {
	PerfTest PerfTests
	Result   PerfResult
}

// These two structures are borrowed from kubernetes/perf-tests:
// https://github.com/kubernetes/perf-tests/blob/master/clusterloader2/pkg/measurement/util/perftype.go
// this is done in order to be compatible with perfdash
type DataItem struct {
	// Data is a map from bucket to real data point (e.g. "Perc90" -> 23.5). Notice
	// that all data items with the same label combination should have the same buckets.
	Data map[string]float64 `json:"data"`
	// Unit is the data unit. Notice that all data items with the same label combination
	// should have the same unit.
	Unit string `json:"unit"`
	// Labels is the labels of the data item.
	Labels map[string]string `json:"labels,omitempty"`
}

// PerfData contains all data items generated in current test.
type PerfData struct {
	// Version is the version of the metrics. The metrics consumer could use the version
	// to detect metrics version change and decide what version to support.
	Version   string     `json:"version"`
	DataItems []DataItem `json:"dataItems"`
	// Labels is the labels of the dataset.
	Labels map[string]string `json:"labels,omitempty"`
}

func getLabelsForTest(summary PerfSummary, metric string) map[string]string {
	node := "other-node"
	if summary.PerfTest.SameNode {
		node = "same-node"
	}
	return map[string]string{
		"scenario":  summary.PerfTest.Scenario,
		"node":      node,
		"test_type": summary.PerfTest.Tool + "-" + summary.PerfTest.Test,
		"metric":    metric,
	}
}

func ExportPerfSummaries(summaries []PerfSummary, reporitDir string) {
	perfData := []DataItem{}
	for _, summary := range summaries {
		if summary.Result.Latency != nil {
			labels := getLabelsForTest(summary, "Latency")
			perfData = append(perfData, summary.Result.Latency.ToPerfData(labels))
		}
		if summary.Result.TransactionRateMetric != nil {
			labels := getLabelsForTest(summary, "TransactionRate")
			perfData = append(perfData, summary.Result.TransactionRateMetric.ToPerfData(labels))
		}
		if summary.Result.ThroughputMetric != nil {
			labels := getLabelsForTest(summary, "Throughput")
			perfData = append(perfData, summary.Result.ThroughputMetric.ToPerfData(labels))
		}
	}
	exportSummary(PerfData{Version: "v1", DataItems: perfData}, reporitDir)
}

func exportSummary(content PerfData, reportDir string) error {
	// this filename needs to be in a specific format for perfdash
	fileName := strings.Join([]string{"NetworkPerformance_benchmark", time.Now().Format(time.RFC3339)}, "_")
	filePath := path.Join(reportDir, strings.Join([]string{fileName, "json"}, "."))
	contentStr, err := prettyPrintJSON(content)
	if err != nil {
		return fmt.Errorf("error formatting summary: %v error: %v", content, err)
	}
	if err := os.WriteFile(filePath, []byte(contentStr), 0644); err != nil {
		return fmt.Errorf("writing to file %v error: %v", filePath, err)
	}
	return nil
}

func prettyPrintJSON(data interface{}) (string, error) {
	output := &bytes.Buffer{}
	if err := json.NewEncoder(output).Encode(data); err != nil {
		return "", fmt.Errorf("building encoder error: %v", err)
	}
	formatted := &bytes.Buffer{}
	if err := json.Indent(formatted, output.Bytes(), "", "  "); err != nil {
		return "", fmt.Errorf("indenting error: %v", err)
	}
	return formatted.String(), nil
}
