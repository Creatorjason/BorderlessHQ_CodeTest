package user_service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/runtime/protoiface"
	"go.mongodb.org/mongo-driver/bson"
	pb "main.go/proto/users"
)

var (
	chars                 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// I have implemented this function to generate a new set of users, based on specified
// generation limit value
func generateNewUsers(numOfUsers int32) []*pb.User {
	allUsers := make([]*pb.User, 0, numOfUsers)
	for i := 0; i < int(numOfUsers); i++ {
		// create individual user
		user := &pb.User{
			Id:        int32(i + 1),
			FirstName: generateRandomWord(7),
			LastName:  generateRandomWord(7),
		}
		allUsers = append(allUsers, user)
	}
	return allUsers

}

func generateRandomWord(ln int) string {
	res := make([]byte, ln)
	for i := range res {
		res[i] = chars[seededRand.Intn(len(chars))]
	}
	return string(res)
}

func protoToJson(m protoiface.MessageV1) []byte {
	marshaler := jsonpb.Marshaler{}
	jsonData, err := marshaler.MarshalToString(m)
	if err != nil {
		fmt.Println("failed to marshal proto to JSON:", err)
		return nil
	}
	return []byte(jsonData)
}
func convertUsersToUserList(usersList []Users) (*pb.UserList, error) {
	userList := &pb.UserList{}

	for _, users := range usersList {
		userWithWallet := &pb.UserWithWallet{
			User: &pb.User{
				Id:        users.User.ID,
				FirstName: users.User.FirstName,
				LastName:  users.User.LastName,
			},
			Wallet: &pb.Wallet{
				Id:       users.Wallet.ID,
				Balance:  users.Wallet.Balance,
				Currency: users.Wallet.Currency,
				UserId:   users.Wallet.UserID,
			},
		}
		userList.UsersWithWallet = append(userList.UsersWithWallet, userWithWallet)
	}

	return userList, nil
}

func getAllWallets(conn *mongo.Client) ([]Wallet, error) {
	var documents []Wallet

	collection := conn.Database("borderlessHQ_service").Collection("borderless_wallets")

	// Get all documents in the collection
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, fmt.Errorf("oops, you don't seem to have any data stored yet in the database, try seeding the database:%v", err)
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &documents)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

// getAllDocuments retrieves all documents from a MongoDB collection
func getAllUsers(conn *mongo.Client) ([]User, error) {
	var documents []User

	collection := conn.Database("borderlessHQ_service").Collection("borderless_users")

	// retrieve all documents in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("oops, you don't seem to have any data stored yet in the database, try seeding the database:%v", err)
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &documents)
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func handleErr(err error) {
	if err != nil {
		log.Println(err)
	}
}
