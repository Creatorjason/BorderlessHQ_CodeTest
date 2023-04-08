package wallet_service

// gRPC client

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "main.go/proto/users"
)

const (
	PORT    = "localhost:9091"
	ConnURI = "mongodb://localhost:27017"
)

type User struct {
	ID        int32  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Wallet    Wallet `json:"wallet"`
}

type Wallet struct {
	ID       int32  `json:"id"`
	Balance  int32  `json:"balance"`
	Currency string `json:"currency"`
	UserID   int32  `json:"userId"`
}

func RunClient() pb.UserServiceClient {

	conn, err := grpc.Dial(PORT, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("unable to connect to grpc server%v ", err.Error())
	}
	// noParams := &pb.NoParams{}

	client := pb.NewUserServiceClient(conn)
	fmt.Println("Connected to server")
	return client
}

func CallSeedUsers(client pb.UserServiceClient, event *pb.SeedUsersEvent, proc chan bool) {
	_, err := client.SeedUsers(context.Background(), event)
	if err != nil {
		log.Println("unable to call SeedUsers from grpc server")
		proc <- false
	}
	fmt.Println("calling SeedUsers()")
	proc <- true
}

func CallListUsers(client pb.UserServiceClient, req *pb.NoParams, proc chan bool) {

	list, err := client.ListUsers(context.Background(), req)

	if err != nil {
		log.Printf("unable to call ListUsers from grpc server: %v", err.Error())
		proc <- false
	}
	jsonp := JsonOutput(list)
	fmt.Println("calling ListUsers()")

	if string(jsonp) == "null" {
		proc <- false
		fmt.Println(" 🔴 oops, you currently do not have any data stored in the database, try seeding the database by choosing option <1>")
		os.Exit(0)
	} else {
		proc <- true
		fmt.Println(string(jsonp))
	}

}
