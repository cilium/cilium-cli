// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package install

import (
	"reflect"
	"testing"

	"helm.sh/helm/v3/pkg/chartutil"
)

func Test_operatorDeploymentValues(t *testing.T) {
	type args struct {
		mergedUserHelmValues chartutil.Values
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "no user input",
			args: args{
				mergedUserHelmValues: map[string]interface{}{},
			},
			want: map[string]string{
				"operator.replicas": "1",
				"operator.updateStrategy.rollingUpdate.maxUnavailable": "100%",
			},
			wantErr: false,
		},
		{
			name: "user provided replica 1",
			args: args{
				mergedUserHelmValues: map[string]interface{}{
					"operator": map[string]interface{}{
						"replicas": "1",
					},
				},
			},
			want: map[string]string{
				"operator.updateStrategy.rollingUpdate.maxUnavailable": "100%",
			},
			wantErr: false,
		},
		{
			name: "user provided replica 2",
			args: args{
				mergedUserHelmValues: map[string]interface{}{
					"operator": map[string]interface{}{
						"replicas": "2",
					},
				},
			},
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name: "user provided update strategy",
			args: args{
				mergedUserHelmValues: map[string]interface{}{
					"operator": map[string]interface{}{
						"updateStrategy": nil,
					},
				},
			},
			want: map[string]string{
				"operator.replicas": "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := operatorDeploymentValues(tt.args.mergedUserHelmValues)
			if (err != nil) != tt.wantErr {
				t.Errorf("operatorDeploymentValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("operatorDeploymentValues() got = %v, want %v", got, tt.want)
			}
		})
	}
}
