package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/andrerampanelli1/fc2-grpc-go/pb/pb"
)

// type UserServiceClient interface {
// 	AddUser(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
// 	AddUserVerbose(ctx context.Context, in *User, opts ...grpc.CallOption) (UserService_AddUserVerboseClient, error)
// 	AddUsers(ctx context.Context, opts ...grpc.CallOption) (UserService_AddUsersClient, error)
// }

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {

	// Insert DB
	fmt.Println(req.Name)

	return &pb.User{
		Id:    "1",
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})

	time.Sleep((time.Second * 3))

	stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User:   &pb.User{},
	})

	time.Sleep((time.Second * 3))

	stream.Send(&pb.UserResultStream{
		Status: "User has been inserted",
		User: &pb.User{
			Id:    "1",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	time.Sleep((time.Second * 3))

	stream.Send(&pb.UserResultStream{
		Status: "User has been inserted",
		User: &pb.User{
			Id:    "1",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}
		if err != nil {
			log.Fatalf("ERROR RECEIVING STREAM %v", err)
		}

		users = append(users, &pb.User{
			Id:    req.GetId(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		})
		fmt.Println("ADDING ", req.GetName())
	}
}

func (*UserService) AddUsersVerbose(stream pb.UserService_AddUsersVerboseServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("ERROR RECEIVING STREAM %v", err)
		}

		fmt.Println("Added ", req)

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   req,
		})

		if err != nil {
			log.Fatalf("ERROR SENDING STREAM %v", err)
		}

	}
}
