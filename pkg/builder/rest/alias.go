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

package rest

import "k8s.io/apiserver/pkg/registry/rest"

// CategoriesProvider if implemented will publish categories for the resource in the Kubernetes discovery service
type CategoriesProvider = rest.CategoriesProvider

// Creater if implemented will expose PUT endpoints for the resource and publish them in the Kubernetes
// discovery service and OpenAPI.
//
// Required for `kubectl apply`.
type Creater = rest.Creater

// CollectionDeleter if implemented will expose DELETE endpoints for resource collections and publish them in
// the Kubernetes discovery service and OpenAPI.
//
// Required for `kubectl delete --all`
type CollectionDeleter = rest.CollectionDeleter

type Connecter = rest.Connecter

type CreaterUpdater = rest.CreaterUpdater

type Exporter = rest.Exporter

// Getter if implemented will expose GET endpoints for the resource and publish them in the Kubernetes
// discovery service and OpenAPI.
//
// Required for `kubectl apply` and most operators.
type Getter = rest.Getter

type GracefulDeleter = rest.GracefulDeleter

// Lister if implemented will enable listing resources.
//
// Required by `kubectl get` and most operators.
type Lister = rest.Lister

// Patcher if implemented will expose POST and GET endpoints for the resource and publish them in the Kubernetes
// discovery service and OpenAPI.
//
// Required by `kubectl apply` and most operators.
type Patcher = rest.Patcher

type Responder = rest.Responder

type Scoper = rest.Scoper

type ShortNamesProvider = rest.ShortNamesProvider

// TableConvertor if implemented will return tabular data from the GET endpoint when requested.
//
// Required by pretty printing `kubectl get`.
type TableConvertor = rest.TableConvertor

// Updater if implemented will expose POST endpoints for the resource and publish them in the Kubernetes
// discovery service and OpenAPI.
//
// Required by `kubectl apply` and most operators.
type Updater = rest.Updater

// Watcher if implemented will enable watching resources.
//
// Required by most operators.
type Watcher = rest.Watcher

// StandardStorage defines the standard endpoints for resources.
type StandardStorage = rest.StandardStorage

type Redirector = rest.Redirector

type Storage = rest.Storage

type ValidateObjectFunc = rest.ValidateObjectFunc

type ValidateObjectUpdateFunc = rest.ValidateObjectUpdateFunc

var NewDefaultTableConvertor = rest.NewDefaultTableConvertor
