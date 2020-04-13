// (c) Copyright IBM Corporation 2019, 2020. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// U.S. Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.
// Licensed Materials - Property of IBM
//
// Copyright (c) 2020 Red Hat, Inc.


// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1.EndpointConfig":       schema_pkg_apis_multicloud_v1alpha1_EndpointConfig(ref),
		"github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1.EndpointConfigSpec":   schema_pkg_apis_multicloud_v1alpha1_EndpointConfigSpec(ref),
		"github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1.EndpointConfigStatus": schema_pkg_apis_multicloud_v1alpha1_EndpointConfigStatus(ref),
	}
}

func schema_pkg_apis_multicloud_v1alpha1_EndpointConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EndpointConfig is the Schema for the endpointconfigs API",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1.EndpointConfigSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1.EndpointConfigStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1.EndpointConfigSpec", "github.com/open-cluster-management/rcm-controller/pkg/apis/multicloud/v1alpha1.EndpointConfigStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_multicloud_v1alpha1_EndpointConfigSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EndpointConfigSpec defines the desired state of EndpointConfig",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}

func schema_pkg_apis_multicloud_v1alpha1_EndpointConfigStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "EndpointConfigStatus defines the observed state of EndpointConfig",
				Properties:  map[string]spec.Schema{},
			},
		},
		Dependencies: []string{},
	}
}
