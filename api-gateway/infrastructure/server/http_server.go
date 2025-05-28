package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tuandq2112/go-microservices/api-gateway/appconfig"
	"github.com/tuandq2112/go-microservices/shared/logger"
	"github.com/tuandq2112/go-microservices/shared/middlewares"
	"github.com/tuandq2112/go-microservices/shared/proto/types/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HttpServer struct {
	Host   string
	Port   string
	mux    *runtime.ServeMux
	logger logger.Logger
}

func NewHttpServer() *HttpServer {
	// Initialize locale

	return &HttpServer{
		Host: appconfig.Host,
		Port: appconfig.Port,
		mux: runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		),
		logger: logger.GetLogger(),
	}
}

func (s *HttpServer) Start() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwmux := s.mux

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := user.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf(
		"%s:%s",
		appconfig.UserServiceHost,
		appconfig.UserServicePort,
	), opts)
	if err != nil {
		s.logger.Error("Failed to register gateway", zap.Error(err))
	}

	s.logger.Info(fmt.Sprintf("API Gateway starting on %s:%s", appconfig.Host, appconfig.Port))

	// Build middleware chain
	var handler http.Handler = s.mux
	// handler = middlewares.JWTMiddleware(middlewares.WhitelistPaths{
	// 	Get:    appconfig.WHITELIST_METHODS_GET_PATH,
	// 	Post:   appconfig.WHITELIST_METHODS_POST_PATH,
	// 	Put:    appconfig.WHITELIST_METHODS_PUT_PATH,
	// 	Delete: appconfig.WHITELIST_METHODS_DELETE_PATH,
	// }, appconfig.USER_CONTEXT_KEY)(handler)
	handler = middlewares.LoggingMiddleware(handler)
	// handler = middlewares.LocaleMiddleware(s.locale)(handler)

	handler = middlewares.EnableCORS(handler)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", appconfig.Host, appconfig.Port), handler); err != nil {
		s.logger.Error("Failed to start server", zap.Error(err))
	}
	return nil
}
