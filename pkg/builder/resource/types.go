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

	"github.com/pwittrock/apiserver-runtime/pkg/builder/resource/resourcestrategy"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apiserver/pkg/registry/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
)

// Object must be implemented by all resources served by the apiserver.
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
			if _, ok := obj.(resourcestrategy.Defaulter); ok {
				s.AddTypeDefaultingFunc(obj, func(o interface{}) {
					o.(resourcestrategy.Defaulter).Default()
				})
			}
		}
		return nil
	}
}
