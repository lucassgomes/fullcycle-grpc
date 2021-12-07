package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/lucassgomes/fullcycle-grpc/pb/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connection to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	AddUsers(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Lucas",
		Email: "lucas@gmail.com",
	}
	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Lucas",
		Email: "lucas@gmail.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}
	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status:", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "1",
			Name:  "Luis",
			Email: "luis@gmail.com",
		},
		{
			Id:    "2",
			Name:  "Lucas",
			Email: "lucas@gmail.com",
		},
		{
			Id:    "3",
			Name:  "Luan",
			Email: "luan@gmail.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 5)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	fmt.Println(res)
}
