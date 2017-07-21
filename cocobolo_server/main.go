package main

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

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

		messages := make(chan *pb.CallbackResponse)

		// Parse input and make HTTP request in a
		// goroutine. Send response back from the goroutine

		go func(in pb.CallbackRequest) {

			// TODO: Check for backoff time
			// Based on the backoff time, write logic
			// to retry request in case of a non 2XX response
			// from the server

			c := &http.Client{
				Timeout: 15 * time.Second,
			}
			resp, err := c.Get(in.URL)
			if err != nil {
				log.Fatal("Could not make request")
			}

			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				log.Fatal("Could not read body")
			}

			bodyString := string(bodyBytes)

			log.Info(bodyString)

			// Create a full callback response object
			// which can be passed back on the channel

			response := &pb.CallbackResponse{RequestId: in.RequestId, Response: bodyString}

			messages <- response

		}(*in)

		response := <-messages

		if err := stream.Send(response); err != nil {
			log.Fatal("Could not send response. %v", err)
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
		log.Fatal("failed to serve: %v", err)
	}
}
