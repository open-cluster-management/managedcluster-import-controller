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
	"os"
	"testing"

	multicloudv1beta1 "github.com/open-cluster-management/endpoint-operator/pkg/apis/multicloud/v1beta1"
	multicloudv1alpha1 "github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func init() {
	os.Setenv("ENDPOINT_CRD_FILE", "../../build/resources/multicloud_v1beta1_endpoint_crd.yaml")
}

func TestNewOperatorDeployment(t *testing.T) {
	type args struct {
		endpointConfig  *multicloudv1alpha1.EndpointConfig
		imageTagPostfix string
	}
	type expectValues struct {
		imageName          string
		imageTagPostfixEnv string
	}

	tests := []struct {
		name string
		args args
		want expectValues
	}{
		{
			name: "Empty Postfix",
			args: args{
				endpointConfig: &multicloudv1alpha1.EndpointConfig{
					Spec: multicloudv1beta1.EndpointSpec{
						ImageRegistry: "sample-registry/uniquePath",
						Version:       "2.3.0",
					},
				},
				imageTagPostfix: "",
			},
			want: expectValues{"sample-registry/uniquePath/endpoint-operator:2.3.0", ""},
		},
		{
			name: "With Postfix Set",
			args: args{
				endpointConfig: &multicloudv1alpha1.EndpointConfig{
					Spec: multicloudv1beta1.EndpointSpec{
						ImageRegistry: "sample-registry-2/uniquePath-2",
						Version:       "1.2.0",
					},
				},
				imageTagPostfix: "-Unique-Postfix",
			},
			want: expectValues{"sample-registry-2/uniquePath-2/endpoint-operator:1.2.0-Unique-Postfix", "-Unique-Postfix"},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			err := os.Setenv(ImageTagPostfixKey, tt.args.imageTagPostfix)
			if err != nil {
				t.Errorf("Cannot set env %s", ImageTagPostfixKey)
			}
			deployment := newOperatorDeployment(tt.args.endpointConfig)
			assert.Equal(t, deployment.Spec.Template.Spec.Containers[0].Image, tt.want.imageName, "image name should match")
			assert.Equal(t, deployment.Spec.Template.Spec.Containers[0].Env[3].Name, ImageTagPostfixKey)
			assert.Equal(t, deployment.Spec.Template.Spec.Containers[0].Env[3].Value, tt.want.imageTagPostfixEnv, "tag postfix should be passed to env")
		})
	}
}

func TestGenerateEndpointCRD(t *testing.T) {
	_, err := GenerateEndpointCRD()
	if err != nil {
		t.Errorf("Cannot generate endpoint crd: %v", err)
		return
	}
}
