package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/andrerampanelli1/fc2-grpc-go/pb/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("COULD NOT CONNECT TO gRPC SERVER: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	AddUser(client)
	AddUserVerbose(client)
	AddUsers(client)
	AddUsersVerbose(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Andre",
		Email: "andrerampanelli1@gmail.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("COULD NOT MAKE gRPC REQUEST: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Andre",
		Email: "andrerampanelli1@gmail.com",
	}

	resStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("COULD NOT MAKE gRPC REQUEST: %v", err)
	}

	for {
		stream, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("COULD NOT RECEIVE THE MESAGE %v", err)
		}

		fmt.Println("STATUS: ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "0",
			Name:  "Andre",
			Email: "andre@email.com",
		},
		&pb.User{
			Id:    "1",
			Name:  "Felipe",
			Email: "felipe@email.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Ricardo",
			Email: "ricardo@email.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Joselito",
			Email: "joselito@email.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Marlene",
			Email: "marlene@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("ERROR CREATING STREAM %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("ERROR RECEIVING RESPONSE %v", err)
	}

	fmt.Println(res)
}

func AddUsersVerbose(client pb.UserServiceClient) {
	stream, err := client.AddUsersVerbose(context.Background())
	if err != nil {
		log.Fatalf("ERROR CREATING STREAM %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "0",
			Name:  "Andre",
			Email: "andre@email.com",
		},
		&pb.User{
			Id:    "1",
			Name:  "Felipe",
			Email: "felipe@email.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Ricardo",
			Email: "ricardo@email.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "Joselito",
			Email: "joselito@email.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Marlene",
			Email: "marlene@email.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user - ", req.GetName())
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("ERROR RECEIVING DATA %v", err)
				break
			}
			fmt.Printf("RECEBENDO USERS\nSTATUS %v - NAME %v\n\n", res.GetStatus(), res.GetUser().GetName())
		}
		close(wait)
	}()

	<-wait
}
