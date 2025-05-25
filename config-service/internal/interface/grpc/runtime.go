package grpc

import (
	"fmt"
	"net"

	"github.com/tuandq2112/go-microservices/config-service/appconfig"
	configHandler "github.com/tuandq2112/go-microservices/config-service/internal/interface/handler/config"
	"github.com/tuandq2112/go-microservices/shared/discovery"
	"github.com/tuandq2112/go-microservices/shared/logger"
	configpb "github.com/tuandq2112/go-microservices/shared/proto/types/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server          *grpc.Server
	configHandler   *configHandler.ConfigHandler
	logger          logger.Logger
	consulRegistrar *discovery.ConsulRegistrar
	healthServer    *health.Server
}

func NewGrpcServer(configHandler *configHandler.ConfigHandler, consulRegistrar *discovery.ConsulRegistrar) (*GrpcServer, error) {
	server := grpc.NewServer()

	// đăng ký reflection service để hỗ trợ grpcurl, debugging, etc
	reflection.Register(server)

	configpb.RegisterConfigServiceServer(server, configHandler)

	zapLogger := logger.GetLogger()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	return &GrpcServer{
		server:          server,
		configHandler:   configHandler,
		logger:          zapLogger,
		consulRegistrar: consulRegistrar,
		healthServer:    healthServer,
	}, nil
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appconfig.Port))
	if err != nil {
		return err
	}
	s.healthServer.SetServingStatus(appconfig.ServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	s.healthServer.SetServingStatus("grpc.health.v1.Health", grpc_health_v1.HealthCheckResponse_SERVING)

	serviceID := fmt.Sprintf("%s-%d", appconfig.ServiceName, appconfig.Port)
	s.consulRegistrar.RegisterService(discovery.ServiceRegistrationConfig{
		ID:      serviceID,
		Name:    appconfig.ServiceName,
		Address: appconfig.Host,
		Port:    appconfig.Port,
		Tags:    []string{"grpc"},
	})
	s.logger.Info(fmt.Sprintf("Starting gRPC server on port %d", appconfig.Port))

	return s.server.Serve(lis)
}
