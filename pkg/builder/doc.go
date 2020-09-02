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

// Package builder contains functions for building Kubernetes apiservers.
//
// Example:
//
//    	err := builder.APIServer.
//		// Setup types and definitions
//		WithSchemeInstallers(yourGroup.Install).
//		WithOpenAPIDefinitions("your-group", "v0.0.0", openapi.GetOpenAPIDefinitions).
//
//		// SimpleResource -- use standard storage and register a subresource under "your-sub-resource"
//		WithResource(&v1alpha1.SimpleResource{}).
//		WithSubResourceAndHandler(
//			&v1alpha1.SimpleResource{}, // parent resource
//			"your-sub-resource", // subresource path
//			&v1alpha1.SubResource{}, // subresource request
//			yourrest.SubResourceHandler). // subresource handler
//
//		// SpecialResource APIs -- use custom storage implementation which implements the API endpoints
//		WithResourceAndHandler(&v1alpha1.SpecialResource{}, yourrest.SpecialResourceHandler).
//
//		// Start the apiserver
//		Execute()
package builder
