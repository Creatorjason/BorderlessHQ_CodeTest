package user_service

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/protobuf/proto"
	pb "main.go/proto/users"
)

// Custom data I used to form the Users output
type User struct {
	ID        int32  `bson:"id"`
	FirstName string `bson:"firstname"`
	LastName  string `bson:"lastname"`
}

type Wallet struct {
	ID       int32  `bson:"id"`
	Balance  int32  `bson:"balance"`
	Currency string `bson:"currency"`
	UserID   int32  `bson:"userid"`
}

type Users struct {
	User   User   `json:"user"`
	Wallet Wallet `json:"wallet"`
}

func (u UserService) SeedUsers(ctx context.Context, req *pb.SeedUsersEvent) (*pb.NoParams, error) {
	// Nats connection
	natsConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("unable to connect to nats server %v", err.Error())
	}
	req.Data = []*pb.User{
		{
			Id:        1,
			FirstName: "Buzz",
			LastName:  "Lightyear",
		},
		{
			Id:        2,
			FirstName: "Lekard",
			LastName:  "Fej",
		},
		{
			Id:        3,
			FirstName: "Sam",
			LastName:  "Stash",
		},
		{
			Id:        4,
			FirstName: "David",
			LastName:  "Oglay",
		},
		{
			Id:        5,
			FirstName: "Justice",
			LastName:  "Nefe",
		},
	}
	// marshal
	payload, err := proto.Marshal(&pb.SeedUsersEvent{
		Id:   req.Id,
		Name: req.Name,
		Data: req.Data,
	})
	if err != nil {
		log.Printf("failed to marshal content %v", err.Error())
	}
	// Push the users array through stream
	if err = natsConn.Publish("users.generated", protoToJson(req)); err != nil {
		log.Fatalf("failed to publish SeedUser event to stream %v", err.Error())
	}
	fmt.Println("sent payload:"+" "+string(payload))
	return &pb.NoParams{}, nil
}

// This method ListUsers retrieves the list of users from the MongoDB database along with their wallets.)
func (us *UserService) ListUsers(ctx context.Context, req *pb.NoParams) (*pb.UserList, error) {
	// var userWithWallet *pb.UserWithWallet
	// var userList *pb.UserList

	allUsers := Users{}
	listOfAllUsers := make([]Users, 0)
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	handleErr(err)
	users, err := getAllUsers(mongoClient)
	// fmt.Printf("%#v\n", users)
	handleErr(err)
	wallets, err := getAllWallets(mongoClient)
	// fmt.Println(wallets)
	handleErr(err)

	// Iterate through all the documents returned per collection
	for idx, user := range users {
		wal := wallets[idx]
		if user.ID == wal.UserID {
			allUsers.User = user
			allUsers.Wallet = wal
			listOfAllUsers = append(listOfAllUsers, allUsers)
		}
	}
	out, _:=convertUsersToUserList(listOfAllUsers)
	return  out, nil

}