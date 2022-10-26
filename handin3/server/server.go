package main

import (
	"context"
	"fmt"
	h "handin3"
	"log"
	"net"
	t "time"

	pb "github.com/katrinesando/Distributed-Systems/tree/ChittyChat_Handin3/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGetCurrentTimeServer
}

var clock h.Vector

func (s *Server) GetTime(ctx context.Context, in *pb.GetTimeRequest) (*pb.GetTimeReply, error) {
	fmt.Printf("Received GetTime request\n")
	return &pb.GetTimeReply{Reply: t.Now().String()}, nil
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGetCurrentTimeServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
