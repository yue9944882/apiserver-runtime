package main

import (
	"github.com/pwittrock/apiserver-runtime/pkg/builder"
	"github.com/pwittrock/apiserver-runtime/pkg/builder/resource"
	"github.com/pwittrock/apiserver-runtime/pkg/example/v1alpha1"
	"github.com/pwittrock/apiserver-runtime/pkg/example/v1beta1"
)

func main() {
	var _ resource.Object = &v1alpha1.ExampleResource{}
	var _ resource.Object = &v1beta1.ExampleResource{}

	cmd, err := builder.APIServer.
		DisableDelegateAuth().
		// v1alpha1 will be the storage version because it was registered first
		WithResource(&v1alpha1.ExampleResource{}).
		// v1beta1 objects will be converted to v1alpha1 versions before being stored
		WithResource(&v1beta1.ExampleResource{}).
		// OpenAPI definitions are optional for an apiserver, unless you need the openapi
		// functionalities for some cases.
		// WithOpenAPIDefinitions("example", "v0.0.0", openapi.GetOpenAPIDefinitions).
		Build()
	if err != nil {
		panic(err)
	}
	// Call Execute on cmd
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
