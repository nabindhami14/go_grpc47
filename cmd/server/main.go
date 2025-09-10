package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	newsv1 "github.com/nabindhami14/go_grpc47/api/news/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	ingrpc "github.com/nabindhami14/go_grpc47/internal/grpc"
	"github.com/nabindhami14/go_grpc47/internal/memstore"

	"buf.build/go/protovalidate"
	protovalidate_interceptor "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"

	"golang.org/x/sync/errgroup"
)

func main() {

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("validator initialization :%v", err)
	}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(protovalidate_interceptor.UnaryServerInterceptor(validator)),
		grpc.ChainStreamInterceptor(protovalidate_interceptor.StreamServerInterceptor(validator)),
	)

	newsv1.RegisterNewsServiceServer(srv, ingrpc.NewServer(memstore.New()))
	healthSrv := health.NewServer()
	healthv1.RegisterHealthServer(srv, healthSrv)

	// IF DATABASE THEN go checkDatabaseHealth -> 1 second interval -> health server status
	// healthSrv.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)
	// healthSrv.SetServingStatus("service-name", healthv1.HealthCheckResponse_SERVING)

	grp, grpCtx := errgroup.WithContext(context.Background())

	grp.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", err)
			}
		}()

		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			err = fmt.Errorf("failed to list: %v", err)
		}

		if listenErr := srv.Serve(lis); listenErr != nil {
			err = fmt.Errorf("failed to serve: %v", listenErr)
		}

		return err
	})

	grp.Go(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", err)
			}
		}()

		interceptSignals(grpCtx)
		healthSrv.Shutdown()
		return shutdown(grpCtx, srv)
	})

	if err := grp.Wait(); err != nil {
		log.Fatal("server shutdown", err)
	}

}

func interceptSignals(ctx context.Context) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		return
	case sig := <-sigc:
		log.Println("intercepted signal: ", sig.String())
		return
	}
}

func shutdown(ctx context.Context, srv *grpc.Server) (err error) {
	done := make(chan struct{})
	go func() {
		srv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		err = fmt.Errorf("grpc server forcefully shutdown: %v", ctx.Err())
		srv.Stop()
	}
	return err

}
