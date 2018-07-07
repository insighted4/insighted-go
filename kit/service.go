package kit

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// HTTPEndpoint encapsulates everything required to build
// an endpoint.
type HTTPEndpoint struct {
	Middleware gin.HandlerFunc
	Methods    map[string]gin.HandlerFunc
}

// Service is the interface of mixed HTTP/gRPC that can be registered and
// hosted by a server. Services provide hooks for service-wide options
// and middlewares and can be used as a means of dependency injection.
type Service interface {
	// HTTPEndpoints default to using a JSON.
	HTTPEndpoints() map[string]HTTPEndpoint

	// RPCMiddleware is for any service-wide gRPC specific middleware
	// for easy integration with 3rd party grpc.UnaryServerInterceptors like
	// http://godoc.org/cloud.google.com/go/trace#Client.GRPCServerInterceptor
	//
	// If you want to apply multiple RPC middlewares,
	// we recommend using:
	// http://godoc.org/github.com/grpc-ecosystem/go-grpc-middleware#ChainUnaryServer
	RPCMiddleware() grpc.UnaryServerInterceptor

	// RPCServiceDesc allows services to declare an alternate gRPC
	// representation of themselves to be hosted on the RPC_PORT (8081 by default).
	RPCServiceDesc() *grpc.ServiceDesc

	// RPCOptions are for service-wide gRPC server options.
	//
	// The underlying kit server already uses the one available grpc.UnaryInterceptor
	// grpc.ServerOption so attempting to pass your own in this method will cause a panic
	// at startup. We recommend using RPCMiddleware() to fill this need.
	RPCOptions() []grpc.ServerOption
}

// Shutdowner allows your service to shutdown gracefully when http server stops.
// This may used when service has any background task which needs to be completed gracefully.
type Shutdowner interface {
	Shutdown()
}
