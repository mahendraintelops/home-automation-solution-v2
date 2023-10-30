package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mahendraintelops/home-automation-solution-v2/user-service/config"
	pb "github.com/mahendraintelops/home-automation-solution-v2/user-service/gen/api/v1"
	grpccontrollers "github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/grpc/server/controllers"
	restcontrollers "github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/rest/server/controllers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinhashubham95/go-actuator"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"

	"github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/rest/client"
)

var (
	host = "localhost"
	port = "25501"
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

	go func() {
		// Start serving gRPC requests
		if err = grpcServer.Serve(lis); err != nil {
			log.Errorf("error serving gRPC: %v", err)
			os.Exit(1)
		}
	}()

	// rest server configuration
	router := gin.Default()
	var restTraceProvider *sdktrace.TracerProvider
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// add opentel
		restTraceProvider = config.InitRestTracer(serviceName, collectorURL, insecure)
		router.Use(otelgin.Middleware(serviceName))
	}
	defer func() {
		if restTraceProvider != nil {
			if err := restTraceProvider.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}
	}()
	// add actuator
	addActuator(router)
	// add prometheus
	addPrometheus(router)

	userController, err := restcontrollers.NewUserController()
	if err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}

	v1 := router.Group("/v1")
	{

		v1.POST("/users", userController.CreateUser)

		v1.GET("/users", userController.ListUsers)

	}

	Port := ":3581"
	log.Println("Server started")
	if err = router.Run(Port); err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}

	// this will not be called as the control won't reach here.
	// call external services here if the HasRestClients value is true
	// (that means this repo is a client to external service(s)
	var err0 error

	bNode14, err0 := client.ExecuteNode14()
	if err0 != nil {
		log.Printf("error occurred: %v", err0)
		return
	}
	log.Printf("response received: %s", string(bNode14))

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func addPrometheus(router *gin.Engine) {
	router.GET("/metrics", prometheusHandler())
}

func addActuator(router *gin.Engine) {
	actuatorHandler := actuator.GetActuatorHandler(&actuator.Config{Endpoints: []int{
		actuator.Env,
		actuator.Info,
		actuator.Metrics,
		actuator.Ping,
		// actuator.Shutdown,
		actuator.ThreadDump,
	},
		Env:     "dev",
		Name:    "user-service",
		Port:    3581,
		Version: "0.0.1",
	})
	ginActuatorHandler := func(ctx *gin.Context) {
		actuatorHandler(ctx.Writer, ctx.Request)
	}
	router.GET("/actuator/*endpoint", ginActuatorHandler)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}
