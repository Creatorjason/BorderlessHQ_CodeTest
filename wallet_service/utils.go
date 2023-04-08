package wallet_service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/runtime/protoiface"
	pb "main.go/proto/users"
)

func ProtoToJson(m protoiface.MessageV1) string {
	marshaler := jsonpb.Marshaler{}
	jsonData, err := marshaler.MarshalToString(m)
	if err != nil {
		fmt.Println("failed to marshal proto to JSON:", err)
		return ""
	}
	return jsonData
}
func JsonOutput(userList *pb.UserList) (jsonOut []byte) {
	var out []User
	for _, ul := range userList.UsersWithWallet {
		wallet := Wallet{
			ID:       ul.Wallet.Id,
			Balance:  ul.Wallet.Balance,
			Currency: ul.Wallet.Currency,
			UserID:   ul.Wallet.UserId,
		}
		user := User{
			ID:        ul.User.Id,
			FirstName: ul.User.FirstName,
			LastName:  ul.User.LastName,
			Wallet:    wallet,
		}
		out = append(out, user)
	}
	jsonOut, err := json.MarshalIndent(out, "", "\t")
	if err != nil{
		log.Panic(err)
	}
	return
}
