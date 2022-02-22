package main

import (
	"context"
	"log"
	"net"

	pb "gRPC-usermngm/proto"

	"database/sql"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50050"
)

// A contructor for the usermanagement server type
func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

//Implementation of the grpc service
type UserManagementServer struct {
	conn *sql.DB
	pb.UnimplementedUserManagementServer
}

//Receiver function of the User Management Type
func (server *UserManagementServer) Run() error {
	//Begin listening on the port specified
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	//register the server as a new grpc service
	pb.RegisterUserManagementServer(srv, server)
	log.Printf("server listening at %v", lis.Addr())
	// Register reflection service on gRPC server.
	reflection.Register(srv)
	return srv.Serve(lis) //return the server
}

// When CreateNewUser function is called we should append any new users to the user management server user_list
// When user is added, read full userlist from file into userlist struct, then append new user and write new userlist back to file
func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetName())
	user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	comand, err := s.conn.Prepare("insert into users(name, age) values ($1,$2)")
	if err != nil {
		log.Println("Unable to insert data:", err.Error())
	}
	comand.Exec(user.Name, user.Age)

	return user, nil
}

func (s *UserManagementServer) DeleteUser(ctx context.Context, u *pb.DelUser) (*pb.UserID, error) {
	coman, err := s.conn.Prepare("delete from users where id=$1")
	if err != nil {
		log.Println("Unable to delete user:", err.Error())
		return nil, err
	}
	coman.Exec(u.GetId())
	log.Printf("Deleted user: %v", u.GetId())

	return &pb.UserID{Id: u.GetId()}, nil
}

func (s *UserManagementServer) GetUser(ctx context.Context, u *pb.UserID) (*pb.User, error) {

	var user *pb.User = &pb.User{}
	comand, err := s.conn.Query("select * from users where id=$1", u.GetId())
	if err != nil {
		log.Println("Error:", err.Error())
	}
	for comand.Next() {
		err = comand.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			log.Println("Error:", err.Error())
		}
	}
	return user, nil
}

// New Reciver function
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var usersList *pb.UserList = &pb.UserList{}
	rows, err := s.conn.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		usersList.Users = append(usersList.Users, &user)
	}
	return usersList, nil
}

func main() {
	//Get a connection to the postgres database
	conn := PgConnect()
	defer conn.Close()

	//First instatiate a new User Management Server
	var server *UserManagementServer = NewUserManagementServer()
	server.conn = conn
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
