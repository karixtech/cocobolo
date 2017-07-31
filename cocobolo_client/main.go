package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "github.com/tsudot/cocobolo/cocobolo"
)

const (
	address = "localhost:50051"
)

func main() {
	// setup connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewCocoboloClient(conn)

	requests := []*pb.CallbackRequest{
		{Endpoint: "https://google.com", Method: "GET", RequestId: "1"},
		{Endpoint: "https://reddit.com", Method: "GET", RequestId: "2"},
	}

	stream, err := client.MakeRequest(context.Background())

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()

			if err == io.EOF {
				close(waitc)
				return
			}

			if err != nil {
				log.Fatalf("Failed to receive message: %v", err)
			}

			log.Printf("Received message. Request id: %s Response: %s", in.RequestId, in.Response)
		}
	}()

	for _, request := range requests {
		if err := stream.Send(request); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
	}

	stream.CloseSend()

	<-waitc
}
