package main

import (
	"context"
	"fmt"
	"github.com/pedrofelli/fc2-grpc/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	//AddUser(client)
	//AddUserVerbose(client)
	//AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient){
	req := &pb.User{
		Id: "0",
		Name: "Pedro",
		Email: "j@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatal("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}


//recive msg with stream of string
func AddUserVerbose(client pb.UserServiceClient){
	req := &pb.User{
		Id: "0",
		Name: "Pedro",
		Email: "j@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatal("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Fatal("Could not receive the msg: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())
	}
}


//send users with sleep
func AddUsers(client pb.UserServiceClient){
	reqs := []*pb.User{
		&pb.User{
			Id: "p1",
			Name: "Pedro",
			Email: "pedro@email.com",
		},
		&pb.User{
			Id: "l1",
			Name: "Lucas",
			Email: "lucas@email.com",
		},
		&pb.User{
			Id: "v1",
			Name: "Vitor",
			Email: "vitor@email.com",
		},
		&pb.User{
			Id: "n1",
			Name: "Nathalia",
			Email: "natalia@email.com",
		},
		&pb.User{
			Id: "j1",
			Name: "Jefferson",
			Email: "jefferson@email.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil{
		log.Fatal("Error creating request: %v", err)
	}

	for _, req := range reqs{
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil{
		log.Fatal("Error reveiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient){
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil{
		log.Fatal("Error create request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id: "p1",
			Name: "Pedro",
			Email: "pedro@email.com",
		},
		&pb.User{
			Id: "l1",
			Name: "Lucas",
			Email: "lucas@email.com",
		},
		&pb.User{
			Id: "v1",
			Name: "Vitor",
			Email: "vitor@email.com",
		},
		&pb.User{
			Id: "n1",
			Name: "Nathalia",
			Email: "natalia@email.com",
		},
		&pb.User{
			Id: "j1",
			Name: "Jefferson",
			Email: "jefferson@email.com",
		},
	}

	wait := make(chan int)

	go func(){
		for _, req := range reqs {
			fmt.Println("Sending user:" , req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()


	go func(){
		for {
			res, err := stream.Recv()
			if err == io.EOF{
				break
			}
			if err != nil {
				log.Fatal("Error receiving data: %", err)
			}
			fmt.Printf("Recebendo user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}

		close(wait)
	}()

	<-wait

}