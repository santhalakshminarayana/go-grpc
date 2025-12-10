// Package interceptors defines interceptor handlers for gRPC server requests
package interceptors

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func RequestRegister(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Printf("calling request: %v", info.FullMethod)
	reqID := uuid.New().String()
	rCtx := context.WithValue(ctx, RequestIDKey, reqID)

	return handler(rCtx, req)
}
