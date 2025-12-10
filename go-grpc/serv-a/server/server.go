// Package server implements the app server functionality
package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
)

func StartServer(host string, port int) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGABRT,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer cancel()

	wg := new(sync.WaitGroup)

	// Register grpc services
	grpcServer := grpcRegister(ctx)

	// Run grpc-server with the config
	runServer(ctx, wg, host, port, grpcServer)

	wg.Wait()
}

func runServer(ctx context.Context, wg *sync.WaitGroup, host string, port int, grpcServer *grpc.Server) {
	lis, err := (&net.ListenConfig{}).Listen(
		ctx,
		"tcp",
		fmt.Sprintf("%v:%v", host, port),
	)
	if err != nil {
		log.Printf("failed to listen for: %v", err)
		return
	}
	log.Printf("Server listeing at %v", lis.Addr().String())

	wg.Add(1)

	// Start gRPC server
	go func() {
		defer wg.Done()

		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("Error serving: %v", err)
		}
	}()

	// Wait for context to cancel and later wait for graceful shutdown
	<-ctx.Done()

	grpcServer.GracefulStop()

	log.Print("Server exit due to shutdown")
}
