// Package service implements gRPC server rpc's
package service

import (
	"context"
	"fmt"

	"github.com/go-grpc/go-proto/serv-a/common/protocommon"
	"github.com/go-grpc/go-proto/serv-a/common/protorpc"
	"github.com/go-grpc/go-proto/serv-a/protouser"
	"github.com/go-grpc/serv-a/interceptors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type userService struct {
	protouser.UnimplementedUserServiceServer
	healthService *health.Server
}

func NewUserService(healthService *health.Server) *userService {
	u := &userService{
		healthService: healthService,
	}
	u.init()

	return u
}

func (u *userService) init() {
	// Init service health as "SERVING"
	u.healthService.SetServingStatus(
		protouser.UserService_ServiceDesc.ServiceName,
		healthpb.HealthCheckResponse_SERVING,
	)
}

func (u *userService) GetUser(ctx context.Context, in *protouser.UserRequest) (*protouser.UserResponse, error) {
	fmt.Printf("request: %v\n", in)

	if in.GetId() >= 10 {
		return nil, status.Errorf(codes.NotFound, "User %v not found", in.GetId())
	}
	userResponse := &protouser.UserResponse{}
	userResponse.SetId(in.GetId())
	userResponse.SetName("user-1")
	userResponse.SetCountry(protocommon.Country_COUNTRY_INDIA)
	userResponse.SetStatus(protouser.UserStatus_USER_STATUS_ACTIVE)

	reqInfo := protorpc.Request_builder{
		RequestId: proto.String(ctx.Value(interceptors.RequestIDKey).(string)),
	}.Build()
	userResponse.SetReqInfo(reqInfo)

	return userResponse, nil
}
