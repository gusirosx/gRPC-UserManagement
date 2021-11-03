package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "gRPC-usermngm/usermgmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":50051"
)

// A contructor for the usermanagement server type
func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

//Implementation of the grpc service
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
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

// When CreateNewUser function is called we should append any new users to the user management server user_list
// When user is added, read full userlist from file into userlist struct, then append new user and write new userlist back to file
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	readBytes, err := ioutil.ReadFile("users.json")
	var users_list *pb.UserList = &pb.UserList{}
	var user_id int32 = int32(rand.Intn(100))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	if err != nil {
		if os.IsNotExist(err) {
			log.Print("File not found. Creating a new file")
			users_list.Users = append(users_list.Users, created_user)
			jsonMarshaling(users_list) // Marshal the users_list into json format
			return created_user, nil
		} else {
			log.Fatalln("Error reading file: ", err)
		}
	}
	if err := protojson.Unmarshal(readBytes, users_list); err != nil {
		log.Fatalf("Failed to parse user list: %v", err)
	}
	users_list.Users = append(users_list.Users, created_user)
	jsonMarshaling(users_list) // Marshal the users_list into json format
	return created_user, nil
}

// Marshal the users_list into json format
func jsonMarshaling(users_list *pb.UserList) {
	jsonBytes, err := protojson.Marshal(users_list)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %v", err)
	}
	if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
		log.Fatalf("Failed write to file: %v", err)
	}
}

// New Reciver function
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalf("Failed read from file: %v", err)
	}
	var users_list *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, users_list); err != nil {
		log.Fatalf("Unmarshaling failed: %v", err)
	}
	return users_list, nil
}

func main() {
	//First instatiate a new User Management Server
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
