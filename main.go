package main

import (
	"fmt"
	"log"
	// "os"

	// "sync"

	// "runtime"

	wm "github.com/dixonwille/wmenu/v5"
	pb "main.go/proto/users"
	us "main.go/user_service"
	ws "main.go/wallet_service"
)

func main() {
	// var wg sync.WaitGroup
	start := make(chan bool)
	proc := make(chan bool)

	go us.StartServer(start)
	ok := <-start
	fmt.Println("=>","🚀🚀 Launching...")

	if ok {
		fmt.Println("  🤩 Hello, welcome to BorderlessHQ 🚀🚀")
		menu := wm.NewMenu(" 📌 What would you like to do?")
		menu.Action(func(opts []wm.Opt) error { handleFunc(opts, proc); return nil })
		menu.Option("Perform database seeding using gRPC", 0, true, nil)
		menu.Option("Get all data from database", 1, false, nil)
		if err := menu.Run(); err != nil {
			log.Fatal(err)
		}
	}
	procComplete := <-proc
	if procComplete {
		fmt.Println("Task completed🕗")
		fmt.Println("Thank you for choosing Borderless Delivery 🔥🔥🔥🔥")
	}
	
	
}

func handleFunc(opts []wm.Opt, proc chan bool) {
	newWS := ws.NewWalletService()
	switch opts[0].Value {
	case 0:
		fmt.Println(" 🔷 This action allows you to call the SeedUsers method on the gRPC, that automatically seeds the database with starting data")
		// payload
		blankUsersArr := make([]*pb.User, 0)
		seedUserEvent := &pb.SeedUsersEvent{
			Id:   1,
			Name: "users.generated",
			Data: blankUsersArr,
		}
		fmt.Println(" 🟢Task in progress, calling SeedUsers()")
		cl := ws.RunClient()
		go ws.CallSeedUsers(cl, seedUserEvent, proc)
		newWS.InsertToUserToDB()

	case 1:
		fmt.Println("🔷This action allows you to call the ListUsers method on the gRPC, that returns a json array of the users with their wallets")
		fmt.Println(" 🟢Task in progress, calling ListUsers()")
		noParams := &pb.NoParams{}
		cl := ws.RunClient()
		go ws.CallListUsers(cl, noParams, proc)
	}
}
