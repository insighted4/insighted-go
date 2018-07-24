package kit

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/insighted4/insighted-go/kit/pprof"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Server encapsulates all logic for registering and running a server.
type Server struct {
	config  Config
	service Service
	logger  logrus.FieldLogger

	httpServer *http.Server
	grpcServer *grpc.Server

	// exit chan for graceful shutdown
	exit chan chan error
}

// New will create a new server for the given Service.
//
// Generally, users should only use the 'Run' function to start a server and use this
// function within tests so they may call ServeHTTP.
func New(cfg Config, svc Service) *Server {
	logger := NewLogger(cfg.LoggerLevel, cfg.LoggerFormat)

	s := &Server{
		config:  cfg,
		service: svc,
		logger:  logger.WithField("component", "server"),
		exit:    make(chan chan error),
	}

	s.httpServer = createHTTPServer(cfg, svc, logger)
	s.grpcServer = createGRPCServer(cfg, svc, logger)

	return s
}

func createGRPCServer(cfg Config, svc Service, logger logrus.FieldLogger) *grpc.Server {
	gdesc := svc.RPCServiceDesc()
	if gdesc == nil {
		return nil
	}

	var inters []grpc.UnaryServerInterceptor
	if mw := svc.RPCMiddleware(); mw != nil {
		inters = append(inters, mw)
	}

	chain := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(inters...))
	interceptors := append(svc.RPCOptions(), chain)

	server := grpc.NewServer(interceptors...)
	server.RegisterService(gdesc, svc)

	return server
}

func createHTTPServer(cfg Config, svc Service, logger logrus.FieldLogger) *http.Server {
	handler := gin.New()
	handler.Use(gin.Recovery())
	handler.Use(CORSHandler())
	handler.Use(LoggerHandler(logger, time.RFC3339, true))
	handler.Use(RequestIDHandler())

	for rel, endpoint := range svc.HTTPEndpoints() {
		group := handler.Group(rel)
		if endpoint.Middleware != nil {
			group.Use(endpoint.Middleware)
		}

		for m, f := range endpoint.Methods {
			group.Handle(m, "", f)
		}
	}

	if cfg.EnablePProf {
		pprof.Register(handler)
	}

	return &http.Server{
		Handler:        handler,
		Addr:           fmt.Sprintf(":%d", cfg.HTTPPort),
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
	}
}

func (s *Server) start() error {
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("HTTP server error - initiating shutting down: %v", err)
			s.stop()
		}
		s.logger.Infof("Listening and serving HTTP on %s\n", s.httpServer.Addr)
	}()

	if s.grpcServer != nil {
		addr := fmt.Sprintf(":%d", s.config.RPCPort)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return errors.Wrap(err, "failed to listen to RPC port")
		}

		go func() {
			err := s.grpcServer.Serve(lis)
			if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
				err = nil
			}

			if err != nil {
				s.logger.Errorf("gRPC server error - initiating shutting down: %v", err)
				s.stop()
			}
		}()

		s.logger.Infof("listening on RPC port: %d", s.config.RPCPort)
	}

	go func() {
		exit := <-s.exit

		// stop listener with timeout
		ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		// stop service
		if shutdown, ok := s.service.(Shutdowner); ok {
			shutdown.Shutdown()
		}

		// stop gRPC server
		if s.grpcServer != nil {
			s.grpcServer.GracefulStop()
		}

		// stop HTTP server
		exit <- s.httpServer.Shutdown(ctx)
	}()

	return nil
}

func (s *Server) stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

// Run will create a new server and register the given
// Service and start up the server(s).
// This will block until the server shuts down.
func Run(cfg Config, svc Service) error {
	srv := New(cfg, svc)

	if err := srv.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	srv.logger.Info("received signal", <-ch)
	return srv.stop()
}
