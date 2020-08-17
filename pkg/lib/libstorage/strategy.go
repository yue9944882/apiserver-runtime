package libstorage

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

// NewStrategy creates and returns a DefaultStrategy instance
func NewStrategy(typer runtime.ObjectTyper) DefaultStrategy {
	return DefaultStrategy{typer, names.SimpleNameGenerator}
}

type ObjectMetaProvider interface {
	GetObjectMeta() *metav1.ObjectMeta
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a ObjectMetaProvider
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	provider, ok := obj.(ObjectMetaProvider)
	if !ok {
		return nil, nil, fmt.Errorf("given object of type %T does not have metadata", obj)
	}
	om := provider.GetObjectMeta()
	return om.GetLabels(), SelectableFields(om), nil
}

// Match is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func Match(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *metav1.ObjectMeta) fields.Set {
	return generic.ObjectMetaFieldsSet(obj, true)
}

type DefaultStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (DefaultStrategy) NamespaceScoped() bool {
	return true
}

func (DefaultStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (DefaultStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (DefaultStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (DefaultStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (DefaultStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (DefaultStrategy) Canonicalize(obj runtime.Object) {
}

func (DefaultStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
