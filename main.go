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

	genericapiserver "k8s.io/apiserver/pkg/server"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	"k8s.io/sample-apiserver/pkg/apis/wardle"
	"k8s.io/sample-apiserver/pkg/apis/wardle/install"
	"k8s.io/sample-apiserver/pkg/apis/wardle/v1alpha1"
	"k8s.io/sample-apiserver/pkg/apis/wardle/v1beta1"
	"k8s.io/sample-apiserver/pkg/apiserver"
	"k8s.io/sample-apiserver/pkg/cmd/server"
	"k8s.io/sample-apiserver/pkg/generated/openapi"
	"k8s.io/sample-apiserver/pkg/lib/librest"
	"k8s.io/sample-apiserver/pkg/lib/libserver"
	"k8s.io/sample-apiserver/pkg/lib/libstorage"
	"k8s.io/sample-apiserver/pkg/registry/wardle/fischer"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	install.Install(apiserver.Scheme)
	librest.AddRESTAPI(v1alpha1.FlunderAPI, libstorage.NewRESTProvider(wardle.FlunderAPI, nil))
	librest.AddRESTAPI(v1alpha1.FischerAPI, libstorage.NewRESTProvider(wardle.FischerAPI,
		fischer.FischerStrategy{}))
	librest.AddRESTAPI(v1beta1.FlunderAPI, libstorage.NewRESTProvider(wardle.FlunderAPI, nil))
	o := libserver.NewOptions(v1alpha1.SchemeGroupVersion, openapi.GetOpenAPIDefinitions)

	stopCh := genericapiserver.SetupSignalHandler()
	cmd := server.NewCommandStartServer(o, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
