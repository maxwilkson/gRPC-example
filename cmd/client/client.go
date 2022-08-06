package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/maxwilkson/grpc-example/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50021", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "João",
		Email: "joao@email.com",
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
		Name:  "Joao",
		Email: "j@j.com",
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
			log.Fatalf("Could not receive the msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "João",
			Email: "joao@email.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Maria",
			Email: "maria@email.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "José",
			Email: "jose@email.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Antônio",
			Email: "antonio@email.com",
		},
		&pb.User{
			Id:    "5",
			Name:  "Francisca",
			Email: "francisca@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUsersStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "1",
			Name:  "João",
			Email: "joao@email.com",
		},
		&pb.User{
			Id:    "2",
			Name:  "Maria",
			Email: "maria@email.com",
		},
		&pb.User{
			Id:    "3",
			Name:  "José",
			Email: "jose@email.com",
		},
		&pb.User{
			Id:    "4",
			Name:  "Antônio",
			Email: "antonio@email.com",
		},
		&pb.User{
			Id:    "5",
			Name:  "Francisca",
			Email: "francisca@email.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
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
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Recebendo user %v com status: %v\n", res.GetUser(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
