package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	// log "github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/tsudot/cocobolo/cocobolo"
)

const (
	port = ":50051"
)

type cocoboloServer struct{}

func (s *cocoboloServer) MakeRequest(stream pb.Cocobolo_MakeRequestServer) error {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

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
			logger.Info("Fetching", zap.String("endpoint", in.Endpoint))

			c := &http.Client{
				Timeout: 15 * time.Second,
			}
			resp, err := c.Get(in.Endpoint)
			if err != nil {
				logger.Fatal("Could not make request", zap.Error(err))
			}

			bodyBytes, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				logger.Fatal("Could not read body")
			}

			bodyString := string(bodyBytes)

			//logger.Info(bodyString)

			// Create a full callback response object
			// which can be passed back on the channel

			response := &pb.CallbackResponse{RequestId: in.RequestId, Response: bodyString}

			messages <- response

		}(*in)

		response := <-messages

		if err := stream.Send(response); err != nil {
			logger.Fatal("Could not send response", zap.Error(err))
		}
	}
}

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	lis, err := net.Listen("tcp", port)

	if err != nil {
		logger.Fatal("failed to listen:", zap.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterCocoboloServer(s, &cocoboloServer{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve:", zap.Error(err))
	}
}
