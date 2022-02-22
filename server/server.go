package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "gRPC-usermngm/proto"

	"github.com/jackc/pgx/v4"
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
	conn                *pgx.Conn
	first_user_creation bool
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
	createSql := ` create table if not exists users(id SERIAL PRIMARY KEY, name text, age int);`
	_, err := s.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}
	s.first_user_creation = false

	log.Printf("Received: %v", in.GetName())

	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}

	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1,$2)",
		created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())
	return created_user, nil
}

// New Reciver function
func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var users_list *pb.UserList = &pb.UserList{}
	rows, err := s.conn.Query(context.Background(), "select * from users")
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
		users_list.Users = append(users_list.Users, &user)
	}
	return users_list, nil
}

func main() {
	//First instatiate a new User Management Server
	password := os.Getenv("PG_PASS")
	user := os.Getenv("PG_USER")
	dbName := os.Getenv("PG_DB_PG")
	host := os.Getenv("PG_HOST")
	connection := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbName)
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	conn, err := pgx.Connect(context.Background(), connection)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())
	user_mgmt_server.conn = conn
	user_mgmt_server.first_user_creation = true
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// _ "github.com/lib/pq"
// )

// // Get the environment credentials and opens a connection with the database
// func DbConect() *sql.DB {
// 	password := os.Getenv("PG_PASS")
// 	user := os.Getenv("PG_USER")
// 	dbName := os.Getenv("PG_DB_STORE")
// 	host := os.Getenv("PG_HOST")
// 	connection := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbName, password, host)
// 	db, err := sql.Open("postgres", connection)
// 	if err != nil {
// 		log.Println("Unable to connect:" + err.Error())
// 	}
// 	return db
// }
