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

package librest

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/server"
	openapicommon "k8s.io/kube-openapi/pkg/common"
	"k8s.io/sample-apiserver/pkg/lib/libtype"
)

var (
	EtcdPath      = "/registry/example.com"
	APIServerName = "example"
	GroupName     = "example.com"
	OpenAPIConfig *openapicommon.Config
	apis          = map[string]map[string]StorageProvider{}
	apiServerFn   func(*server.GenericAPIServer) *server.GenericAPIServer
)

type StorageProvider func(*runtime.Scheme, generic.RESTOptionsGetter) rest.Storage

func AddRESTAPI(api libtype.API, fn StorageProvider) error {
	if _, found := apis[api.Version]; !found {
		apis[api.Version] = map[string]StorageProvider{}
	}
	if _, found := apis[api.Version][api.Resource]; found {
		return errors.Errorf("resource %s %s already registered", api.Version, api.Resource)
	}

	apis[api.Version][api.Resource] = fn
	return nil
}

func InstallAPIs(apiGroupInfo server.APIGroupInfo, s *runtime.Scheme, o generic.RESTOptionsGetter) {
	for version := range apis {
		versionStorage := map[string]rest.Storage{}
		for resource, storage := range apis[version] {
			versionStorage[resource] = storage(s, o)
		}
		apiGroupInfo.VersionedResourcesStorageMap[version] = versionStorage
	}

}

func SetAPIServerFn(fn func(server *server.GenericAPIServer) *server.GenericAPIServer) {
	apiServerFn = fn
}

func GetAPIServer(in *server.GenericAPIServer) *server.GenericAPIServer {
	if apiServerFn != nil {
		return apiServerFn(in)
	}
	return in
}
