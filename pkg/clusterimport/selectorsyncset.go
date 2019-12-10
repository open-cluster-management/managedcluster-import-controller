//Package clusterimport ...
// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package clusterimport

import (
	"context"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const selectorSyncsetName = "multicluster-endpoint"

func selectorSyncsetNsN() types.NamespacedName {
	return types.NamespacedName{
		Name:      selectorSyncsetName,
		Namespace: "",
	}
}

// newSelectorSyncset generate the SelectorSyncset for installing multicluster-endpoint
func newSelectorSyncset() (*hivev1.SelectorSyncSet, error) {
	runtimeObjects, err := generateCommonImportObjects()
	if err != nil {
		return nil, err
	}

	sssNsN := selectorSyncsetNsN()

	runtimeRawExtensions := []runtime.RawExtension{}
	for _, obj := range runtimeObjects {
		runtimeRawExtensions = append(runtimeRawExtensions, runtime.RawExtension{Object: obj})
	}

	return &hivev1.SelectorSyncSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: hivev1.SchemeGroupVersion.String(),
			Kind:       "SelectorSyncSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      sssNsN.Name,
			Namespace: sssNsN.Namespace,
		},
		Spec: hivev1.SelectorSyncSetSpec{
			SyncSetCommonSpec: hivev1.SyncSetCommonSpec{
				Resources: runtimeRawExtensions,
			},
		},
	}, nil
}

// GetSelectorSyncset get the selector syncset use for installing multicluster-endpoint
func GetSelectorSyncset(client client.Client) (*hivev1.SelectorSyncSet, error) {
	sss := &hivev1.SelectorSyncSet{}

	if err := client.Get(context.TODO(), selectorSyncsetNsN(), sss); err != nil {
		return nil, err
	}

	return sss, nil
}

// CreateSelectorSyncset create the selector syncset use for installing multicluster-endpoint
func CreateSelectorSyncset(client client.Client) (*hivev1.SelectorSyncSet, error) {
	sss, err := newSelectorSyncset()
	if err != nil {
		return nil, err
	}

	if err := client.Create(context.TODO(), sss); err != nil {
		return nil, err
	}

	return sss, nil
}
