package main

import (
	"context"
	"github.com/mahendraintelops/home-automation-solution-v2/user-service/config"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	pb "github.com/mahendraintelops/home-automation-solution-v2/user-service/gen/api/v1"
	grpccontrollers "github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/grpc/server/controllers"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

var (
	host = "localhost"
	port = "25500"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func main() {

	// grpc server configuration
	// Initialize the exporter
	var grpcTraceProvider *sdktrace.TracerProvider
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// add opentel
		grpcTraceProvider = config.InitGrpcTracer(serviceName, collectorURL, insecure)
	}
	defer func() {
		if grpcTraceProvider != nil {
			if err := grpcTraceProvider.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}
	}()

	// Set up the TCP listener
	addr := net.JoinHostPort(host, port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("error while starting TCP listener: %v", err)
		os.Exit(1)
	}

	log.Printf("TCP listener started at port: %s", port)

	// Create a new gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Create the User server
	userServer, err := grpccontrollers.NewUserServer()
	if err != nil {
		log.Errorf("error while creating userServer: %v", err)
		os.Exit(1)
	}
	// Register the User server with the gRPC server
	pb.RegisterUserServiceServer(grpcServer, userServer)

	// Enable reflection for the gRPC server
	reflection.Register(grpcServer)

	// Start serving gRPC requests
	if err := grpcServer.Serve(lis); err != nil {
		log.Errorf("error serving gRPC: %v", err)
		os.Exit(1)
	}

}
