package user_service

// Create a gRPC server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "fmt"
	"log"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	pb "main.go/proto/users"
)

type UserService struct {
	pb.UserServiceServer
}

type ServiceCheck struct {
	Name   string `json:"name"`
	Url    string	 `json:"url"`
	Active bool	 `json:"active"`
}

type ServiceResponse struct {
	Id           string  `json:"id"`
	Service      string		 `json:"service"`
	Dependencies []ServiceCheck   `json:"dependencies"`
	Timestamp    string  `json:"timestamp"`
}

const PORT = "localhost:9091"

func StartServer(coord chan bool){
	checks := ServiceHealthCheck()
	if checks == nil {
		log.Fatal("Dependencies are not ready. Exiting...")
		coord <- false

	} else {
		fmt.Println(string(checks))
		lis, err := net.Listen("tcp", PORT)
		if err != nil {
			log.Fatalf("unable to listen %v", err.Error())
			coord <- false
		}
		// fmt.Printf("Listening on port: %v", PORT)
		userGrpcServer := grpc.NewServer()
		pb.RegisterUserServiceServer(userGrpcServer, &UserService{})
		log.Printf("Server running on port %v", PORT)
		go func() {
			if err := userGrpcServer.Serve(lis); err != nil {
				log.Fatalf("unable to create server %v", err.Error())
			}
		}()
		coord <- true
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
	
		userGrpcServer.GracefulStop()
		log.Println("Server stopped successfully")
		os.Exit(1)


	}

}

func ServiceHealthCheck() []byte {
	// checks dependencies health
	natserv := CheckNatsServer()
	mon := CheckMongoServer()

	checks := []ServiceCheck{natserv, mon}
	// checks nats connection
	if !(natserv.Active) && !(mon.Active) {
		return nil
	} else if natserv.Active && mon.Active {
		sr := ServiceResponse{
			Id:           "01EBP4DP4VECW8PHDJJFNEDVKE",
			Service:      "user-service",
			Dependencies: checks,
			Timestamp:    time.Now().String(),
		}
		out, _ := json.MarshalIndent(sr, "", "\t")
		return out
	}
	return nil
}

func CheckNatsServer() ServiceCheck {
	natsConn, err := nats.Connect(nats.DefaultURL, nats.Timeout(10*time.Second))
	if err != nil {
		return ServiceCheck{}
	}
	if natsConn.IsConnected() {
		log.Println("Connected to NATS server successfully!")
		return ServiceCheck{
			Name:   "nats",
			Url:    natsConn.ConnectedUrl(),
			Active: true,
		}
	}
	log.Println("failed to connect to NATS server")
	return ServiceCheck{}
}

func CheckMongoServer() ServiceCheck {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connURI := "mongodb://localhost:27017"
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connURI))
	if err != nil {
		log.Printf("unable to connect to server: %v", err)
		return ServiceCheck{Active: false}
	}
	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		log.Printf("mongodb server not yet ready :%v", err)
		return ServiceCheck{Active: false}
	}
	log.Println("Connected to MongoDB successfully!")
	return ServiceCheck{
		Name:   "mongo",
		Url:    connURI,
		Active: true,
	}
}
