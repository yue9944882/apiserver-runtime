package libserver

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/endpoints/openapi"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	openapicommon "k8s.io/kube-openapi/pkg/common"
	"k8s.io/sample-apiserver/pkg/apiserver"
	"k8s.io/sample-apiserver/pkg/cmd/server"
	"k8s.io/sample-apiserver/pkg/lib/librest"
)

func NewOptions(gv schema.GroupVersion, oa openapicommon.GetOpenAPIDefinitions) *server.ServerOptions {
	librest.GroupName = gv.Group
	librest.EtcdPath = "/registry/" + gv.Group

	// Configure options
	o := &server.ServerOptions{StdOut: os.Stdout, StdErr: os.Stderr}
	o.RecommendedOptions = genericoptions.NewRecommendedOptions(librest.EtcdPath, apiserver.Codecs.LegacyCodec(gv))
	o.RecommendedOptions.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(
		gv, schema.GroupKind{Group: gv.Group})

	// Setup OpenAPI definitions
	server.SetServerConfigFn(func(serverConfig *genericapiserver.RecommendedConfig) *genericapiserver.RecommendedConfig {
		serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(
			oa, openapi.NewDefinitionNamer(apiserver.Scheme))
		serverConfig.OpenAPIConfig.Info.Title = gv.Group
		serverConfig.OpenAPIConfig.Info.Version = "0.1"
		return serverConfig
	})
	return o
}
