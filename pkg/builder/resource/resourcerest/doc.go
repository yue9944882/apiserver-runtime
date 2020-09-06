// Package resourcerest defines interfaces for resource REST implementations.
//
// If a resource implements these interfaces directly on the object, then the resource itself may be used
// as the request handler, and will be registered as the REST handler by default when
// builder.APIServer.WithResource is called.
//
// Alternatively, a REST struct may be defined separately from the object and explicitly registered to handle the
// object with builder.APIServer.WithResourceAndHandler.
package resourcerest
