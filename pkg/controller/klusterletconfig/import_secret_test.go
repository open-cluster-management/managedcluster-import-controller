// (c) Copyright IBM Corporation 2019, 2020. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// U.S. Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.
// Licensed Materials - Property of IBM
//
// Copyright (c) 2020 Red Hat, Inc.

//Package klusterletconfig ...

package klusterletconfig

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	ocinfrav1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	clusterregistryv1alpha1 "k8s.io/cluster-registry/pkg/apis/clusterregistry/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	klusterletv1beta1 "github.com/open-cluster-management/endpoint-operator/pkg/apis/agent/v1beta1"
	klusterletcfgv1beta1 "github.com/open-cluster-management/rcm-controller/pkg/apis/agent/v1beta1"
	"github.com/open-cluster-management/rcm-controller/pkg/controller/clusterregistry"
)

func init() {
	os.Setenv("KLUSTERLET_CRD_FILE", "../../../build/resources/agent.open-cluster-management.io_v1beta1_klusterlet_crd.yaml")
}

func Test_importSecretNsN(t *testing.T) {
	type args struct {
		klusterletConfig *klusterletcfgv1beta1.KlusterletConfig
	}

	tests := []struct {
		name    string
		args    args
		want    types.NamespacedName
		wantErr bool
	}{
		{
			name:    "nil KlusterletConfig",
			args:    args{},
			want:    types.NamespacedName{},
			wantErr: true,
		},
		{
			name: "empty KlusterletConfig.Spec.ClusterName",
			args: args{
				klusterletConfig: &klusterletcfgv1beta1.KlusterletConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "test",
					},
				},
			},
			want:    types.NamespacedName{},
			wantErr: true,
		},
		{
			name: "empty KlusterletConfig.Spec.ClusterNamespace",
			args: args{
				klusterletConfig: &klusterletcfgv1beta1.KlusterletConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "test",
					},
					Spec: klusterletv1beta1.KlusterletSpec{
						ClusterName: "cluster-name",
					},
				},
			},
			want:    types.NamespacedName{},
			wantErr: true,
		},
		{
			name: "no error",
			args: args{
				klusterletConfig: &klusterletcfgv1beta1.KlusterletConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test",
						Namespace: "test",
					},
					Spec: klusterletv1beta1.KlusterletSpec{
						ClusterName:      "cluster-name",
						ClusterNamespace: "cluster-namespace",
					},
				},
			},
			want: types.NamespacedName{
				Name:      "cluster-name" + importSecretNamePostfix,
				Namespace: "cluster-namespace",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := importSecretNsN(tt.args.klusterletConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("importSecretNsN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("importSecretNsN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getImportSecret(t *testing.T) {
	type args struct {
		client           client.Client
		klusterletConfig *klusterletcfgv1beta1.KlusterletConfig
	}

	testKlusterletConfig := &klusterletcfgv1beta1.KlusterletConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-name",
			Namespace: "cluster-namespace",
		},
		Spec: klusterletv1beta1.KlusterletSpec{
			ClusterName:      "cluster-name",
			ClusterNamespace: "cluster-namespace",
		},
	}

	testSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-name" + importSecretNamePostfix,
			Namespace: "cluster-namespace",
		},
	}

	tests := []struct {
		name    string
		args    args
		want    *corev1.Secret
		wantErr bool
	}{
		{
			name: "nil KlusterletConfig",
			args: args{
				client:           fake.NewFakeClient([]runtime.Object{}...),
				klusterletConfig: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "secret does not exist",
			args: args{
				client:           fake.NewFakeClient([]runtime.Object{}...),
				klusterletConfig: testKlusterletConfig,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "secret does exist",
			args: args{
				client:           fake.NewFakeClient([]runtime.Object{testSecret}...),
				klusterletConfig: testKlusterletConfig,
			},
			want:    testSecret,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getImportSecret(tt.args.client, tt.args.klusterletConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("getImportSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getImportSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newImportSecret(t *testing.T) {
	s := scheme.Scheme
	s.AddKnownTypes(clusterregistryv1alpha1.SchemeGroupVersion, &clusterregistryv1alpha1.Cluster{})
	s.AddKnownTypes(klusterletcfgv1beta1.SchemeGroupVersion, &klusterletcfgv1beta1.KlusterletConfig{})
	s.AddKnownTypes(ocinfrav1.SchemeGroupVersion, &ocinfrav1.Infrastructure{})

	klusterletConfig := &klusterletcfgv1beta1.KlusterletConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-name",
			Namespace: "cluster-namespace",
		},
		Spec: klusterletv1beta1.KlusterletSpec{
			ClusterName:      "cluster-name",
			ClusterNamespace: "cluster-namespace",
		},
	}

	cluster := &clusterregistryv1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-name",
			Namespace: "cluster-namespace",
		},
	}

	infrastructConfig := &ocinfrav1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Status: ocinfrav1.InfrastructureStatus{
			APIServerURL: "https://cluster-name.com:6443",
		},
	}

	serviceAccount, err := clusterregistry.NewBootstrapServiceAccount(cluster)
	if err != nil {
		t.Errorf("fail to initialize bootstrap serviceaccount, error = %v", err)
	}

	tokenSecret, err := serviceAccountTokenSecret(serviceAccount)
	if err != nil {
		t.Errorf("fail to initialize serviceaccount token secret, error = %v", err)
	}

	serviceAccount.Secrets = append(serviceAccount.Secrets, corev1.ObjectReference{
		Name: tokenSecret.Name,
	})

	type args struct {
		client           client.Client
		scheme           *runtime.Scheme
		klusterletConfig *klusterletcfgv1beta1.KlusterletConfig
	}

	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "nil scheme",
			args: args{
				client:           fake.NewFakeClient([]runtime.Object{}...),
				scheme:           nil,
				klusterletConfig: nil,
			},
			wantNil: true,
			wantErr: true,
		},
		{
			name: "nil klusterletConfig",
			args: args{
				client:           fake.NewFakeClientWithScheme(s, []runtime.Object{}...),
				scheme:           s,
				klusterletConfig: nil,
			},
			wantNil: true,
			wantErr: true,
		},
		{
			name: "empty klusterletConfig",
			args: args{
				client:           fake.NewFakeClientWithScheme(s, []runtime.Object{}...),
				scheme:           s,
				klusterletConfig: &klusterletcfgv1beta1.KlusterletConfig{},
			},
			wantNil: true,
			wantErr: true,
		},
		{
			name: "no error",
			args: args{
				client: fake.NewFakeClientWithScheme(s, []runtime.Object{
					klusterletConfig,
					cluster,
					infrastructConfig,
					serviceAccount,
					tokenSecret,
					clusterInfoConfigMap(),
				}...),
				scheme:           s,
				klusterletConfig: klusterletConfig,
			},
			wantNil: false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newImportSecret(tt.args.client, tt.args.klusterletConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("newImportSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("newImportSecret() = %v, want %v", got, tt.wantNil)
				return
			}
			if got != nil {
				if got.Data == nil {
					t.Errorf("import secret data should not be empty")
					return
				}
				if len(got.Data["import.yaml"]) == 0 {
					t.Errorf("import.yaml should not be empty")
					return
				}
				if len(got.Data["klusterlet-crd.yaml"]) == 0 {
					t.Errorf("klusterlet-crd.yaml should not be empty")
					return
				}

			}
		})
	}
}

func Test_createImportSecret(t *testing.T) {
	klusterletConfig := &klusterletcfgv1beta1.KlusterletConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-name",
			Namespace: "cluster-namespace",
		},
		Spec: klusterletv1beta1.KlusterletSpec{
			ClusterName:      "cluster-name",
			ClusterNamespace: "cluster-namespace",
		},
	}

	cluster := &clusterregistryv1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-name",
			Namespace: "cluster-namespace",
		},
	}

	infrastructConfig := &ocinfrav1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
		},
		Status: ocinfrav1.InfrastructureStatus{
			APIServerURL: "https://cluster-name.com:6443",
		},
	}
	serviceAccount, err := clusterregistry.NewBootstrapServiceAccount(cluster)
	if err != nil {
		t.Errorf("fail to initialize bootstrap serviceaccount, error = %v", err)
	}

	tokenSecret, err := serviceAccountTokenSecret(serviceAccount)
	if err != nil {
		t.Errorf("fail to initialize serviceaccount token secret, error = %v", err)
	}

	serviceAccount.Secrets = append(serviceAccount.Secrets, corev1.ObjectReference{
		Name: tokenSecret.Name,
	})

	s := scheme.Scheme
	s.AddKnownTypes(clusterregistryv1alpha1.SchemeGroupVersion, &clusterregistryv1alpha1.Cluster{})
	s.AddKnownTypes(klusterletcfgv1beta1.SchemeGroupVersion, &klusterletcfgv1beta1.KlusterletConfig{})
	s.AddKnownTypes(ocinfrav1.SchemeGroupVersion, &ocinfrav1.Infrastructure{})

	fakeClient := fake.NewFakeClientWithScheme(s,
		klusterletConfig,
		cluster,
		infrastructConfig,
		serviceAccount,
		tokenSecret,
		clusterInfoConfigMap(),
	)

	importSecret, err := newImportSecret(fakeClient, klusterletConfig)
	if err != nil {
		t.Errorf("fail to initialize import secret, error = %v", err)
	}
	importSecret.ObjectMeta.OwnerReferences = []metav1.OwnerReference{{
		APIVersion: "clusterregistry.k8s.io/v1alpha1",
		Kind:       "Cluster",
		Name:       "cluster-name",
		UID:        "",
	}}

	type args struct {
		client           client.Client
		scheme           *runtime.Scheme
		cluster          *clusterregistryv1alpha1.Cluster
		klusterletConfig *klusterletcfgv1beta1.KlusterletConfig
	}

	tests := []struct {
		name    string
		args    args
		want    *corev1.Secret
		wantErr bool
	}{
		{
			name: "no error",
			args: args{
				client:           fakeClient,
				scheme:           s,
				cluster:          cluster,
				klusterletConfig: klusterletConfig,
			},
			want:    importSecret,
			wantErr: false,
		},
		{
			name: "secret already exist",
			args: args{
				client: fake.NewFakeClientWithScheme(s,
					klusterletConfig,
					cluster,
					serviceAccount,
					tokenSecret,
					clusterInfoConfigMap(),
					importSecret,
				),
				scheme:           s,
				cluster:          cluster,
				klusterletConfig: klusterletConfig,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createImportSecret(tt.args.client, tt.args.scheme, tt.args.cluster, tt.args.klusterletConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("createImportSecret() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				tt.want.ObjectMeta.ResourceVersion = got.ObjectMeta.ResourceVersion
				tt.want.ObjectMeta.OwnerReferences[0].Controller = got.ObjectMeta.OwnerReferences[0].Controller
				tt.want.ObjectMeta.OwnerReferences[0].BlockOwnerDeletion = got.ObjectMeta.OwnerReferences[0].BlockOwnerDeletion
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createImportSecret() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_toYAML(t *testing.T) {
	testCases := []struct {
		Name    string
		Objects []runtime.Object
		Output  []byte
	}{
		{
			Name:    "no objects",
			Objects: []runtime.Object{},
			Output:  nil,
		},
		{
			Name: "configmap",
			Objects: []runtime.Object{
				&apiextensionv1beta1.CustomResourceDefinition{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "apiextensions.k8s.io/v1beta1",
						Kind:       "CustomResourceDefinition",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "klusterlet.agent.open-cluster-management.io",
					},
					Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
						Group: "agent.open-cluster-management.io",
						Names: apiextensionv1beta1.CustomResourceDefinitionNames{
							Kind:     "Klusterlet",
							ListKind: "KlusterletList",
							Plural:   "klusterlets",
							Singular: "klusterlet",
						},
						Scope: "Namespaced",
					},
				},
				&corev1.ConfigMap{
					TypeMeta: metav1.TypeMeta{
						APIVersion: corev1.SchemeGroupVersion.String(),
						Kind:       "ConfigMap",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cm-test",
						Namespace: "test",
					},
				},
				&corev1.ServiceAccount{
					TypeMeta: metav1.TypeMeta{
						APIVersion: corev1.SchemeGroupVersion.String(),
						Kind:       "ServiceAccount",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      "sa-test",
						Namespace: "test",
					},
				},
			},
			Output: []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  creationTimestamp: null
  name: klusterlet.agent.open-cluster-management.io
spec:
  group: agent.open-cluster-management.io
  names:
    kind: Klusterlet
    listKind: KlusterletList
    plural: klusterlets
    singular: klusterlet
  scope: Namespaced

---
apiVersion: v1
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: cm-test
  namespace: test

---
apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: null
  name: sa-test
  namespace: test
`),
		},
	}
	for _, testCase := range testCases {
		yaml, err := toYAML(testCase.Objects)
		assert.NoError(t, err)
		if !bytes.Equal(testCase.Output, yaml) {
			t.Errorf("toYAML Failed: want %v\n, get %v\n %d %d ", testCase.Output, yaml, len(testCase.Output), len(yaml))
		}
	}
}

func serviceAccountTokenSecret(serviceAccount *corev1.ServiceAccount) (*corev1.Secret, error) {
	if serviceAccount == nil {
		return nil, fmt.Errorf("serviceAccount can not be nil")
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccount.GetName(),
			Namespace: serviceAccount.GetNamespace(),
		},
		Data: map[string][]byte{
			"token": []byte("fake-token"),
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}, nil
}

func clusterInfoConfigMap() *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ibmcloud-cluster-info",
			Namespace: "kube-public",
		},
		Data: map[string]string{
			"cluster_kube_apiserver_host": "api.test.com",
			"cluster_kube_apiserver_port": "6443",
		},
	}
}
