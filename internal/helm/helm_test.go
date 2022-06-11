// SPDX-License-Identifier: Apache-2.0
// Copyright 2022 Authors of Cilium

package helm

import (
	"reflect"
	"testing"

	"github.com/blang/semver/v4"
)

func TestValuesToString(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   map[string]interface{}
		out  string
	}{
		{
			name: "str-slice",
			in:   map[string]interface{}{"hubble": map[string]interface{}{"enabled": true, "metrics": map[string]interface{}{"enabled": []interface{}{"dns", "drop", "tcp", "flow", "icmp", "http"}}}},
			out:  "hubble.enabled=true,hubble.metrics.enabled={dns,drop,tcp,flow,icmp,http}",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := valuesToString("", tt.in); got != tt.out {
				t.Errorf("valuesToString: %q (got) != %q (expected)", got, tt.out)
			}
		})
	}
}

func TestResolveHelmChartVersion(t *testing.T) {
	type args struct {
		versionFlag        string
		chartDirectoryFlag string
	}
	tests := []struct {
		name    string
		args    args
		want    semver.Version
		wantErr bool
	}{
		{
			name:    "valid-version",
			args:    args{versionFlag: "v1.11.5", chartDirectoryFlag: ""},
			want:    semver.Version{Major: 1, Minor: 11, Patch: 5},
			wantErr: false,
		},
		{
			name:    "missing-version",
			args:    args{versionFlag: "v0.0.0", chartDirectoryFlag: ""},
			wantErr: true,
		},
		{
			name:    "invalid-version",
			args:    args{versionFlag: "random-version", chartDirectoryFlag: ""},
			wantErr: true,
		},
		{
			name:    "valid-chart-directory",
			args:    args{versionFlag: "", chartDirectoryFlag: "./testdata"},
			want:    semver.Version{Major: 1, Minor: 2, Patch: 3},
			wantErr: false,
		},
		{
			name:    "invalid-chart-directory",
			args:    args{versionFlag: "", chartDirectoryFlag: "/invalid/chart-directory"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveHelmChartVersion(tt.args.versionFlag, tt.args.chartDirectoryFlag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveHelmChartVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResolveHelmChartVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
