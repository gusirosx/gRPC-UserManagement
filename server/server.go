package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "gRPC-usermngm/usermgmt"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

//Implementation of the grpc service
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

//Defining the service method that is described in the proto file
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100))
	return &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}, nil
}

func main() {
	//Begin listening on the port specified
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//Variable 's' is invoking the new server function from the grpc module
	s := grpc.NewServer()
	//After initialize this new server we're going to register the server as a new grpc service
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("server listening at %v", lis.Addr())
	//Call the server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
