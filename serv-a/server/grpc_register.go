package server

import (
	"context"

	"github.com/go-grpc/go-proto/serv-a/protouser"
	"github.com/go-grpc/serv-a/config"
	"github.com/go-grpc/serv-a/interceptors"
	"github.com/go-grpc/serv-a/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func grpcRegister(ctx context.Context) *grpc.Server {
	config := config.GetConfig()

	// grpc ServerOptions with Interceptors
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptors.RequestRegister),
	}

	// New gRPC server with server options
	grpcServer := grpc.NewServer(opts...)

	// Health service
	healthService := health.NewServer()

	// Initialize services

	// UserService
	userService := service.NewUserService(healthService)
	protouser.RegisterUserServiceServer(grpcServer, userService)

	// Register health service with grpcServer
	healthpb.RegisterHealthServer(grpcServer, healthService)
	// Set health of the server as "SERVING"
	healthService.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// For debugging through grcpurl
	if config.GrpcReflectionEnable {
		reflection.Register(grpcServer)
	}

	return grpcServer
}
