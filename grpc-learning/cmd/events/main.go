package main

import (
	"context"
	"io"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	commonpb "github.com/iypetrov/grpc-learning/gen/common"
	eventspb "github.com/iypetrov/grpc-learning/gen/events"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventsServer struct {
	eventspb.UnimplementedEventsServer
	mu     sync.Mutex
	events []*commonpb.Event
}

func (s *EventsServer) CreateEvent(ctx context.Context, req *eventspb.CreateEventRequest) (*commonpb.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	event := &commonpb.Event{
		Id:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
	}
	s.events = append(s.events, event)
	return event, nil
}

func (s *EventsServer) StreamUpcomingEvents(_ *emptypb.Empty, stream eventspb.Events_StreamUpcomingEventsServer) error {
	for _, e := range s.events {
		if err := stream.Send(e); err != nil {
			return err
		}
	}
	return nil
}

func (s *EventsServer) LiveTicketUpdates(stream eventspb.Events_LiveTicketUpdatesServer) error {
	for {
		update, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Received ticket update: %+v", update)

		if err := stream.Send(update); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	eventspb.RegisterEventsServer(grpcServer, &EventsServer{})
	log.Println("Server is running on port 50051")
	grpcServer.Serve(lis)
}
