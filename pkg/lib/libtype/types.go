package libtype

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type API struct {
	New     func() runtime.Object
	NewList func() runtime.Object
	schema.GroupVersionResource
}
