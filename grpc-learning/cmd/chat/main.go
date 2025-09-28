package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/iypetrov/grpc-learning/gen/chat"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("[UNARY] -> SayHello request received: Name=%q\n", req.Name)
	resp := &pb.HelloResponse{Greeting: "Hello, " + req.Name}
	log.Printf("[UNARY] <- Responding with: %q\n", resp.Greeting)
	return resp, nil
}

func (s *server) StreamMessages(req *pb.StreamRequest, stream pb.ChatService_StreamMessagesServer) error {
	log.Printf("[SERVER-STREAM] -> Start streaming for topic=%q\n", req.Topic)
	for i := 1; i <= 5; i++ {
		msg := fmt.Sprintf("[%s] message %d", req.Topic, i)
		log.Printf("[SERVER-STREAM] -> Sending: %s\n", msg)
		if err := stream.Send(&pb.StreamResponse{Message: msg}); err != nil {
			log.Printf("[SERVER-STREAM] !! Send error: %v", err)
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	log.Println("[SERVER-STREAM] <- Completed streaming")
	return nil
}

func (s *server) UploadMessages(stream pb.ChatService_UploadMessagesServer) error {
	var count int32
	log.Println("[CLIENT-STREAM] -> Waiting for uploaded messages")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[CLIENT-STREAM] <- All messages received. Total=%d\n", count)
			return stream.SendAndClose(&pb.UploadSummary{Count: count})
		}
		if err != nil {
			log.Printf("[CLIENT-STREAM] !! Recv error: %v", err)
			return err
		}
		count++
		log.Printf("[CLIENT-STREAM] -> Received #%d: %q\n", count, req.Content)
	}
}

func (s *server) Chat(stream pb.ChatService_ChatServer) error {
	log.Println("[BIDI] -> Chat session started")
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("[BIDI] <- Client closed stream")
			return nil
		}
		if err != nil {
			log.Printf("[BIDI] !! Recv error: %v", err)
			return err
		}
		log.Printf("[BIDI] -> Received from %s: %q\n", msg.From, msg.Text)

		reply := &pb.ChatMessage{
			From: "Server",
			Text: "Echo: " + msg.Text,
		}
		log.Printf("[BIDI] <- Sending: %q\n", reply.Text)
		if err := stream.Send(reply); err != nil {
			log.Printf("[BIDI] !! Send error: %v", err)
			return err
		}
	}
}

func main() {
	addr := ":50051"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[SERVER] !! Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &server{})

	log.Printf("[SERVER] Listening on %s\n", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[SERVER] !! Serve error: %v", err)
	}
}
