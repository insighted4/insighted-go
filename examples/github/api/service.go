package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/insighted4/insighted-go/kit"
	"google.golang.org/grpc"
)

type Config struct{}

const prefix = "/api/v1"

type service struct {
	client *github.Client
}

var _ GithubProxyServer = service{}

func New() kit.Service {
	return service{
		github.NewClient(nil),
	}
}

func (s service) HTTPEndpoints() map[string]kit.HTTPEndpoint {
	return map[string]kit.HTTPEndpoint{
		"/": {
			Methods: map[string]gin.HandlerFunc{
				http.MethodGet: s.RootHandler,
			},
		},
		prefix + "/users/:id": {
			Middleware: kit.NoCacheHandler(),
			Methods: map[string]gin.HandlerFunc{
				http.MethodGet: s.getUserHTTPHandler,
			},
		},
	}
}

func (s service) RPCMiddleware() grpc.UnaryServerInterceptor {
	return nil
}

func (s service) RPCOptions() []grpc.ServerOption {
	return nil
}

func (s service) RPCServiceDesc() *grpc.ServiceDesc {
	return &_GithubProxy_serviceDesc
}

func (s service) RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Github Service",
	})
}
