/*
Copyright 2017 The Kubernetes Authors.

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

package rest

import (
	"fmt"
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/pwittrock/apiserver-runtime/pkg/apiserver"
	"github.com/pwittrock/apiserver-runtime/pkg/builder/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

func GroupVersionResource(obj runtime.Object) (schema.GroupVersionResource, error) {
	gvks, _, err := apiserver.Scheme.ObjectKinds(obj)
	if err != nil {
		return schema.GroupVersionResource{}, fmt.Errorf(
			"no GroupVersionKind found for %T -- must register the type with the apiserver.Scheme", obj)
	}
	var gvk *schema.GroupVersionKind
	for i := range gvks {
		if gvks[i].Version == runtime.APIVersionInternal {
			continue
		}
		gvk = &gvks[i]
	}
	if gvk == nil {
		return schema.GroupVersionResource{}, fmt.Errorf("no external GroupVersionKind found for %T", obj)
	}

	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: inflection.Plural(strings.ToLower(gvk.Kind)),
	}, nil
}

type HandlerProvider = apiserver.StorageProvider

func New(obj resource.Object) HandlerProvider {
	return func(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (rest.Storage, error) {
		gvr, err := GroupVersionResource(obj)
		if err != nil {
			return nil, err
		}
		s := &DefaultStrategy{
			Object:         obj,
			ObjectTyper:    scheme,
			TableConvertor: rest.NewDefaultTableConvertor(gvr.GroupResource()),
		}
		return newStore(obj.New, obj.NewList, gvr, s, optsGetter)
	}
}

func NewStatus(obj resource.StatusGetSetter) (
	parent resource.Object,
	path string,
	request resource.Object,
	handler HandlerProvider) {

	return obj, "status", obj, func(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (rest.Storage, error) {
		gvr, err := GroupVersionResource(obj)
		if err != nil {
			return nil, err
		}
		s := &StatusSubResourceStrategy{Strategy: &DefaultStrategy{
			Object:         obj,
			ObjectTyper:    scheme,
			TableConvertor: rest.NewDefaultTableConvertor(gvr.GroupResource()),
		}}
		return newStore(obj.New, obj.NewList, gvr, s, optsGetter)
	}
}

func NewWithStrategy(obj resource.Object, s Strategy) HandlerProvider {
	return func(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (rest.Storage, error) {
		gvr, err := GroupVersionResource(obj)
		if err != nil {
			return nil, err
		}
		return newStore(obj.New, obj.NewList, gvr, s, optsGetter)
	}
}

func NewStatusWithStrategy(obj resource.Object, s Strategy) HandlerProvider {
	return func(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (rest.Storage, error) {
		gvr, err := GroupVersionResource(obj)
		if err != nil {
			return nil, err
		}
		s = &StatusSubResourceStrategy{Strategy: s}
		return newStore(obj.New, obj.NewList, gvr, s, optsGetter)
	}
}

// newStore returns a RESTStorage object that will work against API services.
func newStore(
	single, list func() runtime.Object, gvr schema.GroupVersionResource,
	s Strategy, optsGetter generic.RESTOptionsGetter) (*genericregistry.Store, error) {

	store := &genericregistry.Store{
		NewFunc:                  single,
		NewListFunc:              list,
		PredicateFunc:            s.Match,
		DefaultQualifiedResource: gvr.GroupResource(),
		TableConvertor:           s,
		CreateStrategy:           s,
		UpdateStrategy:           s,
		DeleteStrategy:           s,
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: getAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return store, nil
}

// getAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a ObjectMetaProvider
func getAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	provider, ok := obj.(resource.ObjectMetaProvider)
	if !ok {
		return nil, nil, fmt.Errorf("given object of type %T does not have metadata", obj)
	}
	om := provider.GetObjectMeta()
	return om.GetLabels(), selectableFields(om), nil
}

// SelectableFields returns a field set that represents the object.
func selectableFields(obj *metav1.ObjectMeta) fields.Set {
	return generic.ObjectMetaFieldsSet(obj, true)
}

// SubResourceStorageFn is a function that returns objects required to register a subresource into an apiserver
// path is the subresource path from the parent (e.g. "scale"), parent is the resource the subresource
// is under (e.g. &v1.Deployment{}), request is the subresource request (e.g. &Scale{}), storage is
// the storage implementation that handles the requests.
// A SubResourceStorageFn can be used with builder.APIServer.WithSubResourceAndStorageProvider(fn())
type SubResourceStorageFn func() (path string, parent resource.Object, request resource.Object, storage HandlerProvider)

// ResourceStorageFn is a function that returns the objects required to register a resource into an apiserver.
// request is the resource type (e.g. &v1.Deployment{}), storage is the storage implementation that handles
// the requests.
// A ResourceFn can be used with builder.APIServer.WithResourceAndStorageProvider(fn())
type ResourceStorageFn func() (request resource.Object, storage HandlerProvider)
