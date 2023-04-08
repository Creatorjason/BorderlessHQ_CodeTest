package wallet_service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	// "github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/jsonpb"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "google.golang.org/protobuf/proto"
	pb "main.go/proto/users"
	wpb "main.go/proto/wallet"
)

type WalletService struct {
	// nats subscriber
	natsSub *nats.Conn
	// mongoDB client
	mongoClient *mongo.Client
	// user arr
	users []*pb.User
}

// Create a new Wallet Service
func NewWalletService() *WalletService {
	// Create a nats conn
	natsConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("unable to connect to nat server from client %v", err.Error())
	}

	// Create mongoDB client
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(ConnURI))
	if err != nil {
		log.Fatalf("unable to establish connect with MongoDB Atlas %v", err.Error())
	}

	walletService := &WalletService{
		natsSub:     natsConn,
		mongoClient: mongoClient,
	}

	return walletService
}

func (ws *WalletService) getUsersFromStream(subject string) ([]*pb.User, error) {
	allUsers := &pb.SeedUsersEvent{}
	// doneSig := make(chan bool)

	sub, err := ws.natsSub.Subscribe(subject, func(msg *nats.Msg) {

		err := jsonpb.Unmarshal(bytes.NewReader(msg.Data), allUsers)
		fmt.Println(msg.Data)

		if err != nil {
			fmt.Println("failed to unmarshal proto message:", err)
			return
		}
	})
	// handle nats subscription err
	if err != nil {
		fmt.Println("unable to subscribe to Nats subject:", err)
		return nil, err
	}
	fmt.Println(allUsers.Data)
	time.Sleep(5 * time.Second)
	defer sub.Unsubscribe()
	return allUsers.Data, nil
}

func generateNewWallet(userId int32) *wpb.Wallet{
	wallet := &wpb.Wallet{
		Id:       int32(rand.Intn(10000)),
		Balance:  int32(rand.Intn(50)),
		Currency: "NGN",
		UserId:   userId,
	}
	return wallet
}

//Insert users with wallet into mongoDB
func (ws *WalletService) InsertToUserToDB() {
	pay, err := ws.getUsersFromStream("users.generated")
	if err != nil {
		log.Printf("failed to get users from stream: %v", err)
	}
	fmt.Println(pay)
	for _, user := range pay{
	wallet := generateNewWallet(user.Id)
		_, err = ws.mongoClient.Database("borderlessHQ_service").Collection("borderless_users").InsertOne(context.Background(), user)
	
	if err != nil {
		log.Fatalf("Failed to insert users into MongoDB: %v", err.Error())
		return
	} else {
		fmt.Println("Successfully inserted user into MongoDB")
	} 
	_, err = ws.mongoClient.Database("borderlessHQ_service").Collection("borderless_wallets").InsertOne(context.Background(), wallet)
	
	if err != nil {
		log.Fatalf("Failed to insert users into MongoDB: %v", err.Error())
		return
	} else {
		fmt.Println("Successfully inserted wallet into MongoDB")
	} 
	}
}
