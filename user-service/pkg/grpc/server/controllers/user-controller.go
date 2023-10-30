package controllers

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"

	pb "github.com/mahendraintelops/home-automation-solution-v2/user-service/gen/api/v1"
	"github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/grpc/server/models"
	"github.com/mahendraintelops/home-automation-solution-v2/user-service/pkg/grpc/server/services"
)

type UserServer struct {
	userService *services.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserServer() (*UserServer, error) {
	userService, err := services.NewUserService()
	if err != nil {
		return nil, err
	}
	return &UserServer{
		userService: userService,
	}, nil
}

func (*UserServer) Ping(_ context.Context, _ *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Message: "Server is healthy and working!",
	}, nil
}

func (userServer *UserServer) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := models.User{

		Name: req.User.GetName(),
	}

	if _, err := userServer.userService.CreateUser(&user); err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Message: "User Created Successfully!",
	}, nil
}

func (userServer *UserServer) ListUsers(_ context.Context, _ *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users, err := userServer.userService.ListUsers()

	if err != nil {
		return nil, err
	}

	var userList []*pb.User
	for _, retrievedUser := range users {
		userResponse := &pb.User{
			Id: retrievedUser.ID,

			Name: retrievedUser.Name,
		}
		userList = append(userList, userResponse)
	}

	return &pb.ListUsersResponse{
		User: userList,
	}, nil
}

func (userServer *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := req.GetId()
	retrievedUser, err := userServer.userService.GetUser(id)

	if err != nil {
		return nil, err
	}

	userResponse := &pb.User{
		Id: id,

		Name: retrievedUser.Name,
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(ctx)
		currentSpan.SetAttributes(attribute.String("user.id", user.ID))
	}

	return &pb.GetUserResponse{
		User: userResponse,
	}, nil
}
