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

// A contructor for the usermanagement server type
func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UserList{},
	}
}

//Implementation of the grpc service
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	user_list *pb.UserList //array of users struct pointer
}

//Receiver function of the User Management Type
func (server *UserManagementServer) Run() error {
	//Begin listening on the port specified
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//Variable 's' is invoking the new server function from the grpc module
	s := grpc.NewServer()
	//After initialize this new server we're going to register the server as a new grpc service
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	//Call the server
	return s.Serve(lis)
}

// When CreateNewUser function is called we should append any new users to the user management server user list
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	var user_id int32 = int32(rand.Intn(100))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	s.user_list.Users = append(s.user_list.Users, created_user)
	return created_user, nil
}

// New Reciver function
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return s.user_list, nil
}

func main() {
	//First instatiate a new User Management Server
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
