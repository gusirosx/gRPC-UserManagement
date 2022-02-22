package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "gRPC-usermngm/proto"

	"database/sql"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
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
	s := grpc.NewServer() //Variable 's' is invoking the new server function from the grpc module
	//After initialize this new server we're going to register the server as a new grpc service
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis) //Call the server
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
	//First instatiate a new User Management Server
	var server *UserManagementServer = NewUserManagementServer()

	// password := os.Getenv("PG_PASS")
	// user := os.Getenv("PG_USER")
	// dbName := os.Getenv("PG_DB_PG")
	// host := os.Getenv("PG_HOST")
	// connection := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbName, password, host)
	// conn, err := sql.Open("postgres", connection)
	// if err != nil {
	// 	log.Println("Unable to establish connection:", err.Error())
	// }
	conn := PgConnect()
	defer conn.Close()

	server.conn = conn
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func PgConnect() *sql.DB {
	password := os.Getenv("PG_PASS")
	user := os.Getenv("PG_USER")
	dbName := os.Getenv("PG_DB_PG")
	host := os.Getenv("PG_HOST")
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbName, password, host)
	connection, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Unable to establish connection:", err.Error())
	}
	return connection
}
