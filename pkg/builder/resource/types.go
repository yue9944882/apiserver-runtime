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

package resource

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apiserver/pkg/registry/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Object must be implemented by all resources.
type Object interface {
	// Object allows the apiserver libraries to operate on the Object
	runtime.Object

	// ObjectMetaProvider provides the resources ObjectMeta and is required by the apiserver libraries.
	GetObjectMeta() *metav1.ObjectMeta

	// Scoper is used to qualify the resource as either namespace scoped or non-namespace scoped.
	rest.Scoper

	// New returns a new instance of the resource -- e.g. &v1.Deployment{}
	New() runtime.Object

	// NewList return a new list instance of the resource -- e.g. &v1.DeploymentList{}
	NewList() runtime.Object

	// GetGroupVersionResource returns the GroupVersionResource for this resource.  The resource should
	// be the all lowercase and pluralized kind.s
	GetGroupVersionResource() schema.GroupVersionResource

	// IsInternalVersion returns true if the object is also the internal version -- i.e. is the type defined
	// for the API group an alias to this object.
	IsInternalVersion() bool
}

// ListObject contains a list of a particular Object type
type ListObject interface {
	runtime.Object

	GetListMeta() *metav1.ListMeta
}

type AllowCreateOnUpdater interface {
	AllowCreateOnUpdate() bool
}

type AllowUnconditionalUpdater interface {
	AllowUnconditionalUpdate() bool
}

// Canonicalizer functions are invoked before an object is stored to canonicalize the object's format.
// If Canonicalize is implemented fr a type, it will be invoked before storing an object of that type for
// either a create or update.
//
// Canonicalize is only invoked for the type that is the storage version type.
type Canonicalizer interface {
	// Canonicalize formats the object for storage.  Only applied for the version matching the storage version.
	Canonicalize()
}

// Converter defines functions for converting a version of a resource to / from the internal version.
type Converter interface {
	// ConvertFromInternal converts an internal version of the object to this object's version
	ConvertFromInternal(internal interface{})

	// ConvertToInternal converts this version of the object to an internal version of the object.
	ConvertToInternal() (internal interface{})
}

// Defaulter functions are invoked when deserializing an object.  If Default is implemented for a type, the apiserver
// will use it to perform defaulting for that version.
// Default is invoked if the API invocation is for the resource version matching the object type regardless
//of whether or not it is the storage version type for the API.
type Defaulter interface {
	// Default defaults unset values on the object.  Defaults are specific to the version.
	Default()
}

// PrepareForCreater functions are invoked before an object is stored during creation.  If PrepareForCreate
// is implemented for a type, it will be invoked before creating an object of that type.
//
// PrepareForCreater is only invoked for the type that is the storage version type.
type PrepareForCreater interface {
	PrepareForCreate(ctx context.Context)
}

// PrepareForUpdater functions are invoked before an object is stored during update.  If PrepareForCreate
// is implemented for a type, it will be invoked before updating an object of that type.
//
// PrepareForUpdater is only invoked for the type that is the storage version type.
type PrepareForUpdater interface {
	PrepareForUpdate(ctx context.Context, old runtime.Object)
}

// TableConverter functions are invoked when printing an object from `kubectl get`.
type TableConverter interface {
	ConvertToTable(ctx context.Context, tableOptions runtime.Object) (*metav1.Table, error)
}

// Validater functions are invoked before an object is stored to validate the object during creation.  If Validate
// is implemented for a type, it will be invoked before creating an object of that type.
//
// Validater is only invoked for the type that is the storage version type.
type Validater interface {
	Validate(ctx context.Context) field.ErrorList
}

// ValidateUpdater functions are invoked before an object is stored to validate the object during update.
// If ValidateUpdater is implemented for a type, it will be invoked before updating an object of that type.
//
// ValidateUpdater is only invoked for the type that is the storage version type.
type ValidateUpdater interface {
	ValidateUpdate(ctx context.Context, obj runtime.Object) field.ErrorList
}

// StatusGetSetter defines an interface for getting and setting the status for a resource.
type StatusGetSetter interface {
	Object
	// CopyStatus copies the status from the argument to the object
	CopyStatus(ctx context.Context, from runtime.Object)
	// CopySpec copies the spec from the argument to the object
	CopySpec(ctx context.Context, from runtime.Object)
}

// AddToScheme returns a function to add the Objects to the scheme.
//
// AddToScheme will register the objects returned by New and NewList under the GroupVersion for each object.
// AddToScheme will also register the objects under the "__internal" group version for each object that
// returns true for IsInternalVersion.
// AddToScheme will register the defaulting function if it implements the Defaulter inteface.
func AddToScheme(objs ...Object) func(s *runtime.Scheme) error {
	return func(s *runtime.Scheme) error {
		for i := range objs {
			obj := objs[i]
			s.AddKnownTypes(obj.GetGroupVersionResource().GroupVersion(), obj.New(), obj.NewList())
			if obj.IsInternalVersion() {
				s.AddKnownTypes(schema.GroupVersion{
					Group:   runtime.APIVersionInternal,
					Version: obj.GetGroupVersionResource().Version}, obj.New(), obj.NewList())
			}
			if _, ok := obj.(Defaulter); ok {
				s.AddTypeDefaultingFunc(obj, func(o interface{}) {
					o.(Defaulter).Default()
				})
			}
		}
		return nil
	}
}
