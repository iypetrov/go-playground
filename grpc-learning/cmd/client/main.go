package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/iypetrov/grpc-learning/gen/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "gRPC server address")
	doUnary := flag.Bool("unary", false, "Run SayHello (unary) call")
	doSrvStream := flag.Bool("server-stream", false, "Run StreamMessages (server streaming)")
	doCliStream := flag.Bool("client-stream", false, "Run UploadMessages (client streaming)")
	doBiStream := flag.Bool("bidi", false, "Run Chat (bidirectional streaming)")
	name := flag.String("name", "Alice", "Name for SayHello")
	topic := flag.String("topic", "news", "Topic for StreamMessages")
	flag.Parse()

	if !*doUnary && !*doSrvStream && !*doCliStream && !*doBiStream {
		fmt.Fprintln(os.Stderr, "No action specified. Use --unary, --server-stream, --client-stream, --bidi")
		os.Exit(1)
	}

	log.Printf("[CLIENT] Dialing gRPC server at %s\n", *addr)
	conn, err := grpc.NewClient(
		"passthrough:///"+*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("[CLIENT] Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if *doUnary {
		runUnary(ctx, c, *name)
	}
	if *doSrvStream {
		runServerStream(ctx, c, *topic)
	}
	if *doCliStream {
		runClientStream(ctx, c)
	}
	if *doBiStream {
		runBidiStream(ctx, c, *name)
	}

	log.Println("[CLIENT] Done.")
}

func runUnary(ctx context.Context, c pb.ChatServiceClient, name string) {
	log.Println("[UNARY] -> Sending SayHello request")
	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("[UNARY] Error: %v", err)
	}
	log.Printf("[UNARY] <- Response: %s\n", resp.Greeting)
}

func runServerStream(ctx context.Context, c pb.ChatServiceClient, topic string) {
	log.Printf("[SERVER-STREAM] -> Requesting messages on topic: %q\n", topic)
	stream, err := c.StreamMessages(ctx, &pb.StreamRequest{Topic: topic})
	if err != nil {
		log.Fatalf("[SERVER-STREAM] Stream error: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("[SERVER-STREAM] <- Stream ended")
			break
		}
		if err != nil {
			log.Fatalf("[SERVER-STREAM] Recv error: %v", err)
		}
		log.Printf("[SERVER-STREAM] <- %s\n", msg.Message)
	}
}

func runClientStream(ctx context.Context, c pb.ChatServiceClient) {
	log.Println("[CLIENT-STREAM] -> Starting upload of 3 messages")
	stream, err := c.UploadMessages(ctx)
	if err != nil {
		log.Fatalf("[CLIENT-STREAM] Start error: %v", err)
	}
	for i := 1; i <= 3; i++ {
		content := fmt.Sprintf("msg %d @ %s", i, time.Now().Format(time.Stamp))
		log.Printf("[CLIENT-STREAM] -> Sending: %s", content)
		if err := stream.Send(&pb.UploadRequest{Content: content}); err != nil {
			log.Fatalf("[CLIENT-STREAM] Send error: %v", err)
		}
		time.Sleep(300 * time.Millisecond)
	}
	summary, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("[CLIENT-STREAM] Close error: %v", err)
	}
	log.Printf("[CLIENT-STREAM] <- Server summary: sent %d messages\n", summary.Count)
}

func runBidiStream(ctx context.Context, c pb.ChatServiceClient, name string) {
	log.Println("[BIDI] -> Starting chat session")
	stream, err := c.Chat(ctx)
	if err != nil {
		log.Fatalf("[BIDI] Start error: %v", err)
	}

	go func() {
		for i := 1; i <= 3; i++ {
			msg := &pb.ChatMessage{From: name, Text: fmt.Sprintf("Hello %d @ %s", i, time.Now().Format(time.Stamp))}
			log.Printf("[BIDI] -> Sending: %s", msg.Text)
			if err := stream.Send(msg); err != nil {
				log.Printf("[BIDI] Send error: %v", err)
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
		log.Println("[BIDI] -> Closing send side")
		stream.CloseSend()
	}()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("[BIDI] <- Stream ended")
			break
		}
		if err != nil {
			log.Fatalf("[BIDI] Recv error: %v", err)
		}
		log.Printf("[BIDI] <- %s: %s\n", resp.From, resp.Text)
	}
}
