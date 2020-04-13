// (c) Copyright IBM Corporation 2019, 2020. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// U.S. Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.
// Licensed Materials - Property of IBM
//
// Copyright (c) 2020 Red Hat, Inc.

//Package clusterimport ...
package clusterimport

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	multicloudv1beta1 "github.com/open-cluster-management/endpoint-operator/pkg/apis/multicloud/v1beta1"
	multicloudv1alpha1 "github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1"
)

func Test_bootstrapServiceAccountNsN(t *testing.T) {
	type args struct {
		endpointConfig *multicloudv1alpha1.EndpointConfig
	}

	tests := []struct {
		name    string
		args    args
		want    types.NamespacedName
		wantErr bool
	}{
		{
			name: "nil EndpointConfig",
			args: args{
				endpointConfig: nil,
			},
			want:    types.NamespacedName{},
			wantErr: true,
		},
		{
			name: "empty EndpointConfig",
			args: args{
				endpointConfig: &multicloudv1alpha1.EndpointConfig{},
			},
			want:    types.NamespacedName{},
			wantErr: true,
		},
		{
			name: "good EndpointConfig",
			args: args{
				endpointConfig: &multicloudv1alpha1.EndpointConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "endpointConfig",
						Namespace: "namespace",
					},
					Spec: multicloudv1beta1.EndpointSpec{
						ClusterName: "clustername",
					},
				},
			},
			want: types.NamespacedName{
				Name:      "clustername" + BootstrapServiceAccountNamePostfix,
				Namespace: "namespace",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bootstrapServiceAccountNsN(tt.args.endpointConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("bootstrapServiceAccountNsN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bootstrapServiceAccountNsN() = %v, want %v", got, tt.want)
			}
		})
	}
}
