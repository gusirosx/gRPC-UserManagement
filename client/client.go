package main

import (
	"fmt"
	"log"
	"net/http"

	pb "gRPC-usermngm/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// conn holds an open connection to the ping service.
var conn *grpc.ClientConn

// Function to create a new client (to pass the connection to UserManagement function)
func GetClient() pb.UserManagementClient {
	var err error
	conn, err = Connection()
	if err != nil {
		log.Printf("failed to dial server %s: %v", *serverAddr, err)
	}
	return pb.NewUserManagementClient(conn)
}

// Function to create a new user
func CreateNewUser(ctx *gin.Context) {
	client := GetClient()
	defer conn.Close()
	var user pb.NewUser
	// Call BindJSON to bind the received JSON to user
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}
	response, err := client.CreateNewUser(ctx, &user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error when creating the user": err.Error()})
		log.Fatal("Error when calling GetUser:", err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"response": response})
}

func DeleteUser(ctx *gin.Context) {
	client := GetClient()
	defer conn.Close()
	var user pb.User
	// Call BindJSON to bind the received JSON to user
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println(err.Error())
		return
	}
	// Initialize an empty message as input to the call to get users
	response, err := client.DeleteUser(ctx, &pb.DelUser{Id: user.Id})
	if err != nil {
		log.Fatalf("could not delete the selected user: %v", err)
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "could not delete the selected user:" + err.Error()})
	} else {
		// this call should retur an array of users that are stored within the user management servers
		ctx.JSON(http.StatusOK, gin.H{
			"success": "User successfully deleted",
			"userID":  response.Id,
		})
	}

}

func GetUsers(ctx *gin.Context) {
	client := GetClient()
	defer conn.Close()

	// Initialize an empty message as input to the call to get users
	params := &pb.GetUsersParams{}
	response, err := client.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not retrieve users: %v", err)
	}
	// this call should retur an array of users that are stored within the user management servers
	ctx.JSON(http.StatusOK, gin.H{"response": response})
}

func main() {
	// Set up a http server.
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		fmt.Fprintln(ctx.Writer, "Up and running...")
	})

	router.POST("/users", CreateNewUser)
	router.GET("/users", GetUsers)
	router.DELETE("/test", DeleteUser)

	// Run http server
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
