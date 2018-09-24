package api

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/insighted4/insighted-go/kit"
	extensions "github.com/insighted4/insighted-go/kit/extensions"
	"google.golang.org/grpc"
)

const prefix = "/api/v1"

type service struct {
	client *github.Client
	cfg    kit.Config
	logger logrus.FieldLogger
}

var _ GithubProxyServer = service{}

func New(cfg kit.Config) kit.Service {
	return service{
		client: github.NewClient(nil),
		cfg:    cfg,
		logger: kit.NewLogger(cfg.LoggerLevel, cfg.LoggerFormat),
	}
}

func (s service) Config() kit.Config {
	return s.cfg
}

func (s service) HTTPHandler() http.Handler {
	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(extensions.CORSHandler())
	handler.Use(extensions.LoggerHandler(s.logger, time.RFC3339, true))
	handler.Use(extensions.RequestIDHandler())
	handler.NoRoute(extensions.NotFoundHandler)

	handler.GET("/", s.RootHandler)

	group := handler.Use(extensions.NoCacheHandler())
	group.GET(prefix+"/users/:id", s.getUserHTTPHandler)
	return handler
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
