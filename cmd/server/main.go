package main

import (
	"context"
	"log"
	"net"
	"time"

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

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			start := time.Now()

			log.Printf("unary call made with: +%v", info)
			reponse, err := handler(ctx, req)
			log.Printf("time taken: %s", time.Since(start))
			return reponse, err

		},
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			log.Println("second interceptor")
			return handler(ctx, req)
		}),
		grpc.StreamInterceptor(func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			log.Println("server side streaming interceptor")
			return handler(srv, ss)
		}),
	)

	newsv1.RegisterNewsServiceServer(srv, ingrpc.NewServer(memstore.New()))
	healthSrv := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthSrv)

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
