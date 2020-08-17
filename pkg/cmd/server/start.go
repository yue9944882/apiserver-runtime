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

package server

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"github.com/spf13/cobra"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/features"
	"k8s.io/apiserver/pkg/server"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/sample-apiserver/pkg/apiserver"
)

var (
	tokenFile      = flag.String("token-file", "token", "file to write authentication token to")
	writeTokenFile = flag.Bool("write-token-file", false, "if true, write a bearer token to a file")
)

type Startable interface {
	Start(ch <-chan struct{})
}

// ServerOptions contains state for master/api server
type ServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions

	PostStart Startable
	StdOut    io.Writer
	StdErr    io.Writer
}

// NewCommandStartServer provides a CLI handler for 'start master' command
// with a default ServerOptions.
func NewCommandStartServer(o *ServerOptions, stopCh <-chan struct{}) *cobra.Command {
	cmd := &cobra.Command{
		Short: "Launch a API server",
		Long:  "Launch a API server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)

	return cmd
}

// Validate validates ServerOptions
func (o ServerOptions) Validate(args []string) error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in fields required to have valid data
func (o *ServerOptions) Complete() error {
	// register admission plugins

	if err := completeServerOptions(o); err != nil {
		return err
	}

	return nil
}

// Config returns config for the api server given ServerOptions
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	o.RecommendedOptions.Etcd.StorageConfig.Paging = utilfeature.DefaultFeatureGate.Enabled(features.APIListChunking)

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	serverConfig = getServerConfig(serverConfig)

	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
		ExtraConfig:   apiserver.ExtraConfig{},
	}
	return config, nil
}

// RunServer starts a new APIServer given ServerOptions
func (o *ServerOptions) RunServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}

	server.GenericAPIServer.AddPostStartHookOrDie("start-informers", func(context genericapiserver.PostStartHookContext) error {
		config.GenericConfig.SharedInformerFactory.Start(context.StopCh)
		if o.PostStart != nil {
			o.PostStart.Start(context.StopCh)
		}
		if *writeTokenFile {
			if err := ioutil.WriteFile(*tokenFile, []byte(context.LoopbackClientConfig.BearerToken), 0600); err != nil {
				panic(err)
			}
		}
		return nil
	})

	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}

var (
	etcdPath               string
	apiServerFn            func(*server.GenericAPIServer) *server.GenericAPIServer
	serverConfigFn         func(*server.RecommendedConfig) *server.RecommendedConfig
	completeServerConfigFn func(server *ServerOptions) error
)

func SetEctdPath(path string) {
	etcdPath = path
}

func SetCompleteServerOptionsFn(fn func(server *ServerOptions) error) {
	completeServerConfigFn = fn
}

func completeServerOptions(in *ServerOptions) error {
	if completeServerConfigFn != nil {
		return completeServerConfigFn(in)
	}
	return nil
}

func SetServerConfigFn(fn func(server *server.RecommendedConfig) *server.RecommendedConfig) {
	serverConfigFn = fn
}

func getServerConfig(in *server.RecommendedConfig) *server.RecommendedConfig {
	if serverConfigFn != nil {
		return serverConfigFn(in)
	}
	return in
}
