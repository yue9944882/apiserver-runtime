package libstorage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/sample-apiserver/pkg/lib/libtype"
)

var _ rest.StandardStorage = &REST{}

// REST implements a RESTStorage for API services against etcd
type REST struct {
	*genericregistry.Store
}

// NewRESTOrPanic returns a RESTStorage object that will work against API services.  REST functions may be
// overridden by embedding it in another type.  See the rest package for overridable functions
// e.g. the rest.StandardStorage interface.
func newRESTOrPanic(internalAPI libtype.API, strategy interface{},
	scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) *REST {
	store, err := NewStore(internalAPI, strategy, scheme, optsGetter)
	if err != nil {
		panic(err)
	}
	return &REST{Store: store}
}

// NewREST returns a RESTStorage object that will work against API services.
func NewStore(internalAPI libtype.API, strategy interface{},
	scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*genericregistry.Store, error) {
	if strategy == nil {
		strategy = NewStrategy(scheme)
	}

	store := &genericregistry.Store{
		NewFunc:                  internalAPI.New,
		NewListFunc:              internalAPI.NewList,
		PredicateFunc:            Match,
		DefaultQualifiedResource: internalAPI.GroupResource(),

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(internalAPI.GroupResource()),
	}
	if s, ok := strategy.(rest.RESTCreateStrategy); ok {
		store.CreateStrategy = s
	}
	if s, ok := strategy.(rest.RESTUpdateStrategy); ok {
		store.UpdateStrategy = s
	}
	if s, ok := strategy.(rest.RESTDeleteStrategy); ok {
		store.DeleteStrategy = s
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return store, nil
}

// NewRESTProvider provides a function that will create the default REST storage for an API.
func NewRESTProvider(internalAPI libtype.API, strategy interface{}) func(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) rest.Storage {
	return func(scheme *runtime.Scheme, getter generic.RESTOptionsGetter) rest.Storage {
		return newRESTOrPanic(internalAPI, strategy, scheme, getter)
	}
}
