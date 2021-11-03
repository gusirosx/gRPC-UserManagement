package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "gRPC-usermngm/usermgmt"

	"google.golang.org/grpc"
)

const (
	//Address of the grpc server
	address = "localhost:50051"
)

func main() {
	//dial a connection to grpc server
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//withBlock() means that this function will not return until the connection is made
	defer conn.Close()

	//create a new client (to pass the connection to that function)
	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var new_users = make(map[string]int32)
	new_users["Alice"] = 43
	new_users["Bob"] = 30
	//call the create new user function by looping over the new users map
	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		//r - response from the grpc server
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		// show on terminal the User Details
		log.Printf(`User Details:
		NAME: %s
		AGE: %d
		ID: %d`, r.GetName(), r.GetAge(), r.GetId())
	}
	// Initialize an empty message as input to the call to get users
	params := &pb.GetUsersParams{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not retrieve users: %v", err)
	}
	log.Print("\nUSER LIST: \n")
	// this call should retur an array of users that are stored within the user management servers
	fmt.Printf("r.GetUsers(): %v\n", r.GetUsers())
}
