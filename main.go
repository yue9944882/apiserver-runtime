/*
Copyright 2016 The Kubernetes Authors.

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

package main

import (
	"flag"
	"os"

	"github.com/pwittrock/apiserver-runtime/pkg/apis/wardle/v1alpha1"
	"github.com/pwittrock/apiserver-runtime/pkg/apis/wardle/v1beta1"
	"github.com/pwittrock/apiserver-runtime/pkg/apiserver"
	"github.com/pwittrock/apiserver-runtime/pkg/cmd/server"
	"github.com/pwittrock/apiserver-runtime/pkg/generated/openapi"
	wardleregistry "github.com/pwittrock/apiserver-runtime/pkg/registry"
	fischerstorage "github.com/pwittrock/apiserver-runtime/pkg/registry/wardle/fischer"
	flunderstorage "github.com/pwittrock/apiserver-runtime/pkg/registry/wardle/flunder"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	apiserver.GroupName = "wardle.example.com"
	apiserver.APIs[v1alpha1.SchemeGroupVersion.WithResource("flunders")] = func(s *runtime.Scheme, g generic.RESTOptionsGetter) (rest.Storage, error) {
		return wardleregistry.RESTInPeace(flunderstorage.NewREST(s, g)), nil
	}
	apiserver.APIs[v1alpha1.SchemeGroupVersion.WithResource("fischers")] = func(s *runtime.Scheme, g generic.RESTOptionsGetter) (rest.Storage, error) {
		return wardleregistry.RESTInPeace(fischerstorage.NewREST(s, g)), nil
	}
	apiserver.APIs[v1beta1.SchemeGroupVersion.WithResource("flunders")] = func(s *runtime.Scheme, g generic.RESTOptionsGetter) (rest.Storage, error) {
		return wardleregistry.RESTInPeace(flunderstorage.NewREST(s, g)), nil
	}
	server.SetOpenAPIDefinitions("Wardle", "0.1", openapi.GetOpenAPIDefinitions)

	stopCh := genericapiserver.SetupSignalHandler()
	options := server.NewWardleServerOptions(os.Stdout, os.Stderr, v1alpha1.SchemeGroupVersion)
	cmd := server.NewCommandStartWardleServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
