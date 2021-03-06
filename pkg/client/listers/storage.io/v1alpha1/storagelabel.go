/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/tommenx/storage/pkg/apis/storage.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// StorageLabelLister helps list StorageLabels.
type StorageLabelLister interface {
	// List lists all StorageLabels in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.StorageLabel, err error)
	// StorageLabels returns an object that can list and get StorageLabels.
	StorageLabels(namespace string) StorageLabelNamespaceLister
	StorageLabelListerExpansion
}

// storageLabelLister implements the StorageLabelLister interface.
type storageLabelLister struct {
	indexer cache.Indexer
}

// NewStorageLabelLister returns a new StorageLabelLister.
func NewStorageLabelLister(indexer cache.Indexer) StorageLabelLister {
	return &storageLabelLister{indexer: indexer}
}

// List lists all StorageLabels in the indexer.
func (s *storageLabelLister) List(selector labels.Selector) (ret []*v1alpha1.StorageLabel, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.StorageLabel))
	})
	return ret, err
}

// StorageLabels returns an object that can list and get StorageLabels.
func (s *storageLabelLister) StorageLabels(namespace string) StorageLabelNamespaceLister {
	return storageLabelNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// StorageLabelNamespaceLister helps list and get StorageLabels.
type StorageLabelNamespaceLister interface {
	// List lists all StorageLabels in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.StorageLabel, err error)
	// Get retrieves the StorageLabel from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.StorageLabel, error)
	StorageLabelNamespaceListerExpansion
}

// storageLabelNamespaceLister implements the StorageLabelNamespaceLister
// interface.
type storageLabelNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all StorageLabels in the indexer for a given namespace.
func (s storageLabelNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.StorageLabel, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.StorageLabel))
	})
	return ret, err
}

// Get retrieves the StorageLabel from the indexer for a given namespace and name.
func (s storageLabelNamespaceLister) Get(name string) (*v1alpha1.StorageLabel, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("storagelabel"), name)
	}
	return obj.(*v1alpha1.StorageLabel), nil
}
