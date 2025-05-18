package main

import (
	"log"
	"net"

	newsv1 "github.com/nabindhami14/go_grpc47/api/news/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	ingrpc "github.com/nabindhami14/go_grpc47/internal/grpc"
	"github.com/nabindhami14/go_grpc47/internal/memstore"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	srv := grpc.NewServer()
	newsv1.RegisterNewsServiceServer(srv, ingrpc.NewServer(memstore.New()))
	healthSrv := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthSrv)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
