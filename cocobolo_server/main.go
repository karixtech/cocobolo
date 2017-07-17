package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/tsudot/cocobolo/cocobolo"
)

const (
	port = ":50051"
)

type cocoboloServer struct{}

func (s *cocoboloServer) MakeRequest(stream pb.Cocobolo_MakeRequestServer) error {
	for {
		in, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Println(in.URL)

		messages := make(chan *pb.CallbackResponse)

		// Parse input and make HTTP request in a
		// goroutine. Send response back from the goroutine

		go func(in pb.CallbackRequest) {
			c := &http.Client{
				Timeout: 15 * time.Second,
			}
			resp, err := c.Get(in.URL)
			if err != nil {
				log.Fatalf("Could not make request")
			}

			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				log.Fatalf("Could not read body")
			}

			bodyString := string(bodyBytes)

			log.Printf(bodyString)

			response := &pb.CallbackResponse{RequestId: in.RequestId, Response: bodyString}
			messages <- response

		}(*in)

		response := <-messages

		if err := stream.Send(response); err != nil {
			log.Fatalf("Could not send response. %v", err)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCocoboloServer(s, &cocoboloServer{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
